package tests

import (
	"testing"

	"fin.orion.dev/internal/commitlint"
)

func TestCommitValidator(t *testing.T) {
	validator := commitlint.NewValidator(nil)

	tests := []struct {
		name    string
		message string
		want    bool
		errors  []string
	}{
		{
			name:    "Valid feat commit",
			message: "feat: add new feature",
			want:    true,
		},
		{
			name:    "Valid fix commit with scope",
			message: "fix(auth): correct password validation",
			want:    true,
		},
		{
			name:    "Valid docs commit",
			message: "docs: update README",
			want:    true,
		},
		{
			name:    "Valid test commit with scope",
			message: "test(api): add endpoint tests",
			want:    true,
		},
		{
			name:    "Valid chore commit",
			message: "chore(deps): update dependencies",
			want:    true,
		},
		{
			name:    "Valid breaking change",
			message: "feat!: breaking change",
			want:    true,
		},
		{
			name:    "Valid commit with body",
			message: "feat: add feature\n\nThis is the body of the commit.",
			want:    true,
		},
		{
			name:    "Valid commit with footer",
			message: "fix: correct bug\n\nCloses: #123",
			want:    true,
		},
		{
			name:    "Invalid type",
			message: "invalid: add feature",
			want:    false,
			errors:  []string{"Tipo 'invalid' não é válido"},
		},
		{
			name:    "Missing type",
			message: ": add feature",
			want:    false,
			errors:  []string{"Formato inválido"},
		},
		{
			name:    "Missing description",
			message: "feat:",
			want:    false,
			errors:  []string{"Formato inválido"},
		},
		{
			name:    "Empty message",
			message: "",
			want:    false,
			errors:  []string{"Mensagem de commit vazia"},
		},
		{
			name:    "Description too long",
			message: "feat: " + string(make([]byte, 100)),
			want:    false,
			errors:  []string{"Descrição muito longa"},
		},
		{
			name:    "Description ends with period",
			message: "feat: add feature.",
			want:    false,
			errors:  []string{"Descrição não deve terminar com ponto"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			commit := validator.ValidateCommit(tt.message)

			if commit.IsValid != tt.want {
				t.Errorf("ValidateCommit() = %v, want %v", commit.IsValid, tt.want)
			}

			if !tt.want && len(commit.Errors) == 0 {
				t.Errorf("Expected errors but got none")
			}

			if tt.want && len(commit.Errors) > 0 {
				t.Errorf("Expected no errors but got: %v", commit.Errors)
			}
		})
	}
}

func TestCommitTypes(t *testing.T) {
	types := commitlint.GetCommitTypes()
	expectedTypes := []string{
		"feat", "fix", "docs", "style", "refactor",
		"test", "chore", "perf", "ci", "build", "revert",
	}

	if len(types) != len(expectedTypes) {
		t.Errorf("Expected %d types, got %d", len(expectedTypes), len(types))
	}

	for _, expectedType := range expectedTypes {
		found := false
		for _, actualType := range types {
			if actualType == expectedType {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected type %s not found", expectedType)
		}
	}
}

func TestFormatCommitMessage(t *testing.T) {
	tests := []struct {
		name        string
		commitType  string
		scope       string
		description string
		want        string
	}{
		{
			name:        "Format without scope",
			commitType:  "feat",
			scope:       "",
			description: "add new feature",
			want:        "feat: add new feature",
		},
		{
			name:        "Format with scope",
			commitType:  "fix",
			scope:       "auth",
			description: "correct password validation",
			want:        "fix(auth): correct password validation",
		},
		{
			name:        "Format with complex scope",
			commitType:  "test",
			scope:       "api-endpoints",
			description: "add integration tests",
			want:        "test(api-endpoints): add integration tests",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := commitlint.FormatCommitMessage(tt.commitType, tt.scope, tt.description)
			if got != tt.want {
				t.Errorf("FormatCommitMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCommitMessageParsing(t *testing.T) {
	validator := commitlint.NewValidator(nil)

	tests := []struct {
		name             string
		message          string
		expectedType     commitlint.CommitType
		expectedScope    string
		expectedDesc     string
		expectedBody     string
		expectedFooter   string
		expectedBreaking bool
	}{
		{
			name:             "Simple commit",
			message:          "feat: add new feature",
			expectedType:     commitlint.TypeFeat,
			expectedScope:    "",
			expectedDesc:     "add new feature",
			expectedBody:     "",
			expectedFooter:   "",
			expectedBreaking: false,
		},
		{
			name:             "Commit with scope",
			message:          "fix(auth): correct validation",
			expectedType:     commitlint.TypeFix,
			expectedScope:    "auth",
			expectedDesc:     "correct validation",
			expectedBody:     "",
			expectedFooter:   "",
			expectedBreaking: false,
		},
		{
			name:             "Commit with body",
			message:          "docs: update README\n\nThis is the body.",
			expectedType:     commitlint.TypeDocs,
			expectedScope:    "",
			expectedDesc:     "update README",
			expectedBody:     "This is the body.",
			expectedFooter:   "",
			expectedBreaking: false,
		},
		{
			name:             "Commit with footer",
			message:          "fix: correct bug\n\nCloses: #123",
			expectedType:     commitlint.TypeFix,
			expectedScope:    "",
			expectedDesc:     "correct bug",
			expectedBody:     "",
			expectedFooter:   "Closes: #123",
			expectedBreaking: false,
		},
		{
			name:             "Breaking change with !",
			message:          "feat!: breaking change",
			expectedType:     commitlint.TypeFeat,
			expectedScope:    "",
			expectedDesc:     "breaking change",
			expectedBody:     "",
			expectedFooter:   "",
			expectedBreaking: true,
		},
		{
			name:             "Breaking change in footer",
			message:          "feat: new feature\n\nBREAKING CHANGE: This is breaking",
			expectedType:     commitlint.TypeFeat,
			expectedScope:    "",
			expectedDesc:     "new feature",
			expectedBody:     "",
			expectedFooter:   "BREAKING CHANGE: This is breaking",
			expectedBreaking: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			commit := validator.ValidateCommit(tt.message)

			if commit.Type != tt.expectedType {
				t.Errorf("Type = %v, want %v", commit.Type, tt.expectedType)
			}

			if commit.Scope != tt.expectedScope {
				t.Errorf("Scope = %v, want %v", commit.Scope, tt.expectedScope)
			}

			if commit.Description != tt.expectedDesc {
				t.Errorf("Description = %v, want %v", commit.Description, tt.expectedDesc)
			}

			if commit.Body != tt.expectedBody {
				t.Errorf("Body = %v, want %v", commit.Body, tt.expectedBody)
			}

			if commit.Footer != tt.expectedFooter {
				t.Errorf("Footer = %v, want %v", commit.Footer, tt.expectedFooter)
			}

			if commit.Breaking != tt.expectedBreaking {
				t.Errorf("Breaking = %v, want %v", commit.Breaking, tt.expectedBreaking)
			}
		})
	}
}

func TestValidatorConfig(t *testing.T) {
	// Test custom configuration
	config := &commitlint.Config{
		MaxSubjectLength: 50,
		MinSubjectLength: 5,
		AllowedTypes:     []commitlint.CommitType{commitlint.TypeFeat, commitlint.TypeFix},
		RequireScope:     true,
		RequireBody:      true,
		RequireFooter:    false,
	}

	validator := commitlint.NewValidator(config)

	// Test valid commit with custom config
	validCommit := "feat(scope): valid description\n\nThis is the body."
	commit := validator.ValidateCommit(validCommit)

	if !commit.IsValid {
		t.Errorf("Expected valid commit, got errors: %v", commit.Errors)
	}

	// Test invalid commit (missing scope)
	invalidCommit := "feat: missing scope"
	commit = validator.ValidateCommit(invalidCommit)

	if commit.IsValid {
		t.Errorf("Expected invalid commit due to missing scope")
	}

	// Test invalid commit (missing body)
	invalidCommit2 := "feat(scope): missing body"
	commit = validator.ValidateCommit(invalidCommit2)

	if commit.IsValid {
		t.Errorf("Expected invalid commit due to missing body")
	}

	// Test invalid commit (invalid type)
	invalidCommit3 := "docs(scope): invalid type\n\nbody"
	commit = validator.ValidateCommit(invalidCommit3)

	if commit.IsValid {
		t.Errorf("Expected invalid commit due to invalid type")
	}
}

func TestIsGitRepository(t *testing.T) {
	// This test will depend on the environment
	// We can't easily test this without mocking, so we'll just call it
	_ = commitlint.IsGitRepository()
}

func TestGetCommitTypes(t *testing.T) {
	types := commitlint.GetCommitTypes()

	// Should return all valid types
	if len(types) != len(commitlint.ValidTypes) {
		t.Errorf("Expected %d types, got %d", len(commitlint.ValidTypes), len(types))
	}

	// Should contain all expected types
	expectedTypes := map[string]bool{
		"feat":     true,
		"fix":      true,
		"docs":     true,
		"style":    true,
		"refactor": true,
		"test":     true,
		"chore":    true,
		"perf":     true,
		"ci":       true,
		"build":    true,
		"revert":   true,
	}

	for _, commitType := range types {
		if !expectedTypes[commitType] {
			t.Errorf("Unexpected type: %s", commitType)
		}
	}
}
