package enum

//go:generate go-enum --marshal --values --nocomments --nocase

// AnnotationType
// ENUM
//
//	enum="@Enum"
//
// )

type Annotation string

const AnnotationEnum Annotation = "@enum"
