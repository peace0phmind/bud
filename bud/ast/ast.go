package ast

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

func ParseFromFile(inputFile string) ([]byte, error) {
	fileSet := token.NewFileSet()
	_, err := parser.ParseFile(fileSet, inputFile, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("generate: error parsing input file '%s': %s", inputFile, err)
	}

	return nil, nil
}

func GetRecvType(fd *ast.FuncDecl) *ast.TypeSpec {
	if fd.Recv != nil {
		if fd.Recv.NumFields() == 1 {
			var recvTypeIdent *ast.Ident
			switch tt := fd.Recv.List[0].Type.(type) {
			case *ast.Ident:
				recvTypeIdent = tt

			case *ast.StarExpr:
				if itt, ok := tt.X.(*ast.Ident); ok {
					recvTypeIdent = itt
				}
			}

			if recvTypeIdent != nil && recvTypeIdent.Obj != nil {
				if recvType, ok := recvTypeIdent.Obj.Decl.(*ast.TypeSpec); ok {
					return recvType
				}
			}
		}
	}

	return nil
}

func InspectMapper[From any, To any](fileNode *ast.File, fileSet *token.FileSet, mapper func(*From) *To) []*To {
	result := []*To{}

	ast.Inspect(fileNode, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.TypeSpec:
			if ts, ok := any(x).(*From); ok {
				if x.Doc == nil {
					x.Doc = FindDocLocationCommentGroup(fileNode, fileSet, x.Pos())
				}

				if x.Comment == nil {
					x.Comment = FindCommentLocationCommentGroup(fileNode, fileSet, x.Pos())
				}

				if t := mapper(ts); t != nil {
					result = append(result, t)
				}
			}
		case *ast.FuncDecl:
			if fd, ok := any(x).(*From); ok {
				if x.Doc == nil {
					x.Doc = FindDocLocationCommentGroup(fileNode, fileSet, x.Pos())
				}

				if x.Recv != nil {
					recvType := GetRecvType(x)
					if recvType != nil {
						if recvType.Doc == nil {
							recvType.Doc = FindDocLocationCommentGroup(fileNode, fileSet, recvType.Pos())
						}

						if recvType.Comment == nil {
							recvType.Comment = FindCommentLocationCommentGroup(fileNode, fileSet, recvType.Pos())
						}
					}
				}

				for _, field := range x.Type.Params.List {
					if field.Doc == nil {
						field.Doc = FindDocLocationCommentGroup(fileNode, fileSet, field.Pos())
					}
					if field.Comment == nil {
						field.Comment = FindCommentLocationCommentGroup(fileNode, fileSet, field.Pos())
					}
				}

				if t := mapper(fd); t != nil {
					result = append(result, t)
				}
			}
		}

		return true
	})

	return result
}

func FindDocLocationCommentGroup(fileNode *ast.File, fileSet *token.FileSet, pos token.Pos) *ast.CommentGroup {
	indentPos := fileSet.Position(pos)

	for _, commentGroup := range fileNode.Comments {
		commentGroupPos := fileSet.Position(commentGroup.End())

		if commentGroupPos.Line+1 == indentPos.Line && commentGroupPos.Offset < indentPos.Offset {
			return commentGroup
		}
	}

	return nil
}

func FindCommentLocationCommentGroup(fileNode *ast.File, fileSet *token.FileSet, pos token.Pos) *ast.CommentGroup {
	indentPos := fileSet.Position(pos)

	for _, commentGroup := range fileNode.Comments {
		commentGroupPos := fileSet.Position(commentGroup.End())

		if commentGroupPos.Line == indentPos.Line && indentPos.Offset < commentGroupPos.Offset {
			return commentGroup
		}
	}

	return nil
}
