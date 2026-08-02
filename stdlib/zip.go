package stdlib

import (
	"errors"
	"bytes"
	"io"

	"github.com/2dprototype/tender"
	"github.com/yeka/zip"
)

var zipModule = map[string]tender.Object{
	"writer": &tender.NativeFunction{Name: "writer", Value: zipNewWriter}, 
	"reader": &tender.NativeFunction{Name: "reader", Value: zipNewReader}, 
}

func zipNewWriter(args ...tender.Object) (tender.Object, error) {
	if len(args) != 0 { 
		return nil, tender.ErrWrongNumArguments 
	} 
	
	var zipBuffer bytes.Buffer 
	zipWriter := zip.NewWriter(&zipBuffer) 

	return &tender.ImmutableMap{ 
		Value: map[string]tender.Object{ 
			"create": &tender.NativeFunction{ 
				Value: func(args ...tender.Object) (tender.Object, error) {
					if len(args) != 2 && len(args) != 3 { 
						return nil, tender.ErrWrongNumArguments
					}
					filename, _ := tender.ToString(args[0]) 
					content, _ := tender.ToByteSlice(args[1]) 
					
					var f io.Writer
					var err error

					if len(args) == 3 {
						password, _ := tender.ToString(args[2])
						f, err = zipWriter.Encrypt(filename, password, zip.StandardEncryption)
					} else {
						f, err = zipWriter.Create(filename)
					}
					
					if err != nil { 
						return wrapError(err), nil 
					}

					_, err = f.Write(content) 
					if err != nil { 
						return wrapError(err), nil 
					}

					return nil, nil 
				},
			},
			"bytes": &tender.NativeFunction{ 
				Value: func(args ...tender.Object) (tender.Object, error) { 
					if len(args) != 0 { 
						return nil, tender.ErrWrongNumArguments 
					} 
					return &tender.Bytes{Value: zipBuffer.Bytes()}, nil 
				}, 
			},
			"close": &tender.NativeFunction{ 
				Value: FuncARE(zipWriter.Close),
			},	
			"flush": &tender.NativeFunction{ 
				Value: FuncARE(zipWriter.Flush),
			},
			// "set_comment": &tender.NativeFunction{ 
				// // yeka/zip Writer lacks SetComment. Stubbed to maintain backwards compatibility.
				// Value: func(args ...tender.Object) (tender.Object, error) { 
					// return nil, nil 
				// }, 
			// },
			"set_offset": &tender.NativeFunction{ 
				Value: FuncAI64R(zipWriter.SetOffset),
			},
		},
	}, nil
}

func zipNewReader(args ...tender.Object) (tender.Object, error) {
	if len(args) != 1 { 
		return nil, tender.ErrWrongNumArguments 
	} 
	data, err := ToFileData(args[0])
	if err != nil {
		return nil, err
	}
	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil { 
		return wrapError(err), nil 
	} 
	
	files := make([]tender.Object, len(reader.File)) 
	
	for i, file := range reader.File { 
		files[i] = makeZipFile(file) 
	} 

	return &tender.ImmutableMap{ 
		Value: map[string]tender.Object{ 
			"files" : &tender.Array{Value: files}, 
			"comment" : &tender.String{Value: reader.Comment}, 
		},
	}, nil 
}

func makeZipFile(file *zip.File) *tender.ImmutableMap {
	// Track the password state locally so we don't crash yeka/zip
	hasPassword := false 

	return &tender.ImmutableMap{ 
		Value: map[string]tender.Object{ 
			"name" : &tender.String{Value: file.Name},
			"comment" : &tender.String{Value: file.Comment},
			"non_utf8" : tender.FromBool(file.Flags&0x800 == 0), 
			"creator_version" : &tender.Int{Value: int64(file.CreatorVersion)},
			"reader_version" : &tender.Int{Value: int64(file.ReaderVersion)},
			"flags" : &tender.Int{Value: int64(file.Flags)},
			"method" : &tender.Int{Value: int64(file.Method)},
			"modified" : &tender.Time{Value: file.ModTime()},
			"modified_time" : &tender.Int{Value: int64(file.ModifiedTime)},
			"modified_date" : &tender.Int{Value: int64(file.ModifiedDate)},
			"crc32" : &tender.Int{Value: int64(file.CRC32)},
			"compressed_size" : &tender.Int{Value: int64(file.CompressedSize64)},
			"uncompressed_size" : &tender.Int{Value: int64(file.UncompressedSize64)},
			"extra" : &tender.Bytes{Value: file.Extra},
			"external_attrs" : &tender.Int{Value: int64(file.ExternalAttrs)},
			"defer_auth" : tender.FromBool(file.DeferAuth),
			"is_encrypted": tender.FromBool(file.IsEncrypted()),
			"mode": &tender.Int{Value: int64(file.Mode())},
			"data_offset": &tender.NativeFunction{
				Value: func(args ...tender.Object) (tender.Object, error) {
					offset, err := file.DataOffset()
					if err != nil {
						return wrapError(err), nil
					}
					return &tender.Int{Value: offset}, nil
				},
			},
			"set_password": &tender.NativeFunction{
				Value: func(args ...tender.Object) (tender.Object, error) {
					if len(args) != 1 {
						return nil, tender.ErrWrongNumArguments
					}
					password, _ := tender.ToString(args[0])
					file.SetPassword(password)
					hasPassword = true
					return nil, nil
				},
			},
			"read": &tender.NativeFunction{ 
				Value: func(args ...tender.Object) (tender.Object, error) {
					if len(args) != 0 { 
						return nil, tender.ErrWrongNumArguments
					}
					
					if file.IsEncrypted() && !hasPassword {
						return wrapError(errors.New("file is encrypted: please call set_password() before reading")), nil
					}

					fileReader, err := file.Open()
					if err != nil { 
						return wrapError(err), nil 
					} 
					defer fileReader.Close() 
					
					content, err := io.ReadAll(fileReader) 
					if err != nil { 
						return wrapError(err), nil 
					} 
					
					return &tender.Bytes{Value: content}, nil 
				},
			},
		},
	}
}