package linecode

type Language struct {
	Name      string
	Slug      string
	Extension string
	Comment   Comment
}

func (l Language) String() string {
	return l.Name
}

func FindLanguage(slug string) Language {
	for _, l := range Languages {
		if l.Slug == slug {
			return l
		}
	}
	return Language{}
}

type Comment struct {
	Start string
	End   string
}

var CStyleComment = Comment{
	Start: "/**",
	End:   "**/",
}

var Languages = []Language{
	{
		Name:      "C++",
		Slug:      "cpp",
		Extension: "cpp",
		Comment:   CStyleComment,
	},
	{
		Name:      "Java",
		Slug:      "java",
		Extension: "java",
		Comment:   CStyleComment,
	},
	{
		Name:      "Python",
		Slug:      "python",
		Extension: "py",
		Comment: Comment{
			Start: `"""`,
			End:   `"""`,
		},
	},
	{
		Name:      "Python3",
		Slug:      "python3",
		Extension: "py",
		Comment: Comment{
			Start: `"""`,
			End:   `"""`,
		},
	},

	{
		Name:      "C",
		Slug:      "c",
		Extension: "c",
		Comment:   CStyleComment,
	},
	{
		Name:      "C#",
		Slug:      "csharp",
		Extension: "cs",
		Comment:   CStyleComment,
	},
	{
		Name:      "JavaScript",
		Slug:      "javascript",
		Extension: "js",
		Comment:   CStyleComment,
	},
	{
		Name:      "Ruby",
		Slug:      "ruby",
		Extension: "rb",
		Comment: Comment{
			Start: `=begin`,
			End: `=end	`,
		},
	},
	{
		Name:      "Swift",
		Slug:      "swift",
		Extension: "swift",
		Comment:   CStyleComment,
	},
	{
		Name:      "Go",
		Slug:      "golang",
		Extension: "go",
		Comment:   CStyleComment,
	},
	{
		Name:      "Scala",
		Slug:      "scala",
		Extension: "scala",
		Comment:   CStyleComment,
	},
	{
		Name:      "Kotlin",
		Slug:      "kotlin",
		Extension: "kt",
		Comment:   CStyleComment,
	},
	{
		Name:      "Rust",
		Slug:      "rust",
		Extension: "rs",
		Comment:   CStyleComment,
	},
	{
		Name:      "PHP",
		Slug:      "php",
		Extension: "php",
		Comment:   CStyleComment,
	},
}
