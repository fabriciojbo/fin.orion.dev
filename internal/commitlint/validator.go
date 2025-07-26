package commitlint

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

// CommitType representa os tipos válidos de commit
type CommitType string

const (
	// Tipos de commit válidos
	TypeFeat     CommitType = "feat"
	TypeFix      CommitType = "fix"
	TypeDocs     CommitType = "docs"
	TypeStyle    CommitType = "style"
	TypeRefactor CommitType = "refactor"
	TypeTest     CommitType = "test"
	TypeChore    CommitType = "chore"
	TypePerf     CommitType = "perf"
	TypeCI       CommitType = "ci"
	TypeBuild    CommitType = "build"
	TypeRevert   CommitType = "revert"
)

// ValidTypes lista todos os tipos válidos
var ValidTypes = []CommitType{
	TypeFeat, TypeFix, TypeDocs, TypeStyle, TypeRefactor,
	TypeTest, TypeChore, TypePerf, TypeCI, TypeBuild, TypeRevert,
}

// CommitMessage representa uma mensagem de commit parseada
type CommitMessage struct {
	Type        CommitType
	Scope       string
	Description string
	Body        string
	Footer      string
	Breaking    bool
	IsValid     bool
	Errors      []string
}

// Config representa a configuração do commitlint
type Config struct {
	MaxSubjectLength int
	MinSubjectLength int
	AllowedTypes     []CommitType
	RequireScope     bool
	RequireBody      bool
	RequireFooter    bool
}

// DefaultConfig retorna a configuração padrão
func DefaultConfig() *Config {
	return &Config{
		MaxSubjectLength: 72,
		MinSubjectLength: 1,
		AllowedTypes:     ValidTypes,
		RequireScope:     false,
		RequireBody:      false,
		RequireFooter:    false,
	}
}

// Validator representa o validador de commits
type Validator struct {
	config *Config
}

// NewValidator cria um novo validador
func NewValidator(config *Config) *Validator {
	if config == nil {
		config = DefaultConfig()
	}
	return &Validator{config: config}
}

// ValidateCommit valida uma mensagem de commit
func (v *Validator) ValidateCommit(message string) *CommitMessage {
	commit := &CommitMessage{
		IsValid: true,
		Errors:  []string{},
	}

	// Parse da mensagem
	v.parseMessage(message, commit)

	// Só validar se o parse foi bem-sucedido
	if commit.IsValid {
		v.validateType(commit)
		v.validateScope(commit)
		v.validateDescription(commit)
		v.validateBody(commit)
		v.validateFooter(commit)
	}

	return commit
}

// parseMessage faz o parse da mensagem de commit
func (v *Validator) parseMessage(message string, commit *CommitMessage) {
	lines := strings.Split(message, "\n")
	if len(lines) == 0 || strings.TrimSpace(lines[0]) == "" {
		commit.IsValid = false
		commit.Errors = append(commit.Errors, "Mensagem de commit vazia")
		return
	}

	// Parse da primeira linha (header)
	header := strings.TrimSpace(lines[0])
	v.parseHeader(header, commit)

	// Parse do body e footer
	if len(lines) > 1 {
		v.parseBodyAndFooter(lines[1:], commit)
	}
}

// parseHeader faz o parse do header do commit
func (v *Validator) parseHeader(header string, commit *CommitMessage) {
	// Regex para conventional commits com suporte a breaking change (!)
	// formato: type!:(...) ou type(scope)!: ...
	pattern := `^(\w+)(!)?(?:\(([^)]+)\))?:\s*(.+)$`
	re := regexp.MustCompile(pattern)

	matches := re.FindStringSubmatch(header)
	if len(matches) < 5 {
		commit.IsValid = false
		commit.Errors = append(commit.Errors, "Formato inválido. Use: type(scope): description")
		return
	}

	commit.Type = CommitType(strings.TrimSpace(matches[1]))
	commit.Breaking = matches[2] == "!"
	commit.Scope = strings.TrimSpace(matches[3])
	commit.Description = strings.TrimSpace(matches[4])
}

// parseBodyAndFooter faz o parse do body e footer
func (v *Validator) parseBodyAndFooter(lines []string, commit *CommitMessage) {
	var bodyLines []string
	var footerLines []string
	inFooter := false

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Verificar se é footer (linha que começa com BREAKING CHANGE, Closes, etc.)
		if strings.HasPrefix(line, "BREAKING CHANGE:") ||
			strings.HasPrefix(line, "Closes:") ||
			strings.HasPrefix(line, "Fixes:") ||
			strings.HasPrefix(line, "See also:") {
			inFooter = true
		}

		if inFooter {
			footerLines = append(footerLines, line)
		} else {
			bodyLines = append(bodyLines, line)
		}
	}

	commit.Body = strings.Join(bodyLines, "\n")
	commit.Footer = strings.Join(footerLines, "\n")

	// Verificar breaking change no footer
	if strings.Contains(commit.Footer, "BREAKING CHANGE:") {
		commit.Breaking = true
	}
}

// validateType valida o tipo do commit
func (v *Validator) validateType(commit *CommitMessage) {
	if commit.Type == "" {
		commit.IsValid = false
		commit.Errors = append(commit.Errors, "Formato inválido")
		return
	}

	// Verificar se o tipo é válido
	valid := false
	for _, validType := range v.config.AllowedTypes {
		if commit.Type == validType {
			valid = true
			break
		}
	}

	if !valid {
		commit.IsValid = false
		commit.Errors = append(commit.Errors,
			fmt.Sprintf("Tipo '%s' não é válido", commit.Type))
	}
}

// validateScope valida o escopo do commit
func (v *Validator) validateScope(commit *CommitMessage) {
	if v.config.RequireScope && commit.Scope == "" {
		commit.IsValid = false
		commit.Errors = append(commit.Errors, "Escopo é obrigatório")
	}
}

// validateDescription valida a descrição do commit
func (v *Validator) validateDescription(commit *CommitMessage) {
	if commit.Description == "" {
		commit.IsValid = false
		commit.Errors = append(commit.Errors, "Formato inválido")
		return
	}

	length := len(commit.Description)
	if length < v.config.MinSubjectLength {
		commit.IsValid = false
		commit.Errors = append(commit.Errors, "Descrição muito curta")

	}

	if length > v.config.MaxSubjectLength {
		commit.IsValid = false
		commit.Errors = append(commit.Errors, "Descrição muito longa")
	}

	// Verificar se não termina com ponto
	if strings.HasSuffix(commit.Description, ".") {
		commit.IsValid = false
		commit.Errors = append(commit.Errors, "Descrição não deve terminar com ponto")
	}
}

// validateBody valida o body do commit
func (v *Validator) validateBody(commit *CommitMessage) {
	if v.config.RequireBody && commit.Body == "" {
		commit.IsValid = false
		commit.Errors = append(commit.Errors, "Body é obrigatório")
	}
}

// validateFooter valida o footer do commit
func (v *Validator) validateFooter(commit *CommitMessage) {
	if v.config.RequireFooter && commit.Footer == "" {
		commit.IsValid = false
		commit.Errors = append(commit.Errors, "Footer é obrigatório")
	}
}

// GetCommitMessage obtém a mensagem do último commit
func GetCommitMessage() (string, error) {
	cmd := exec.Command("git", "log", "-1", "--pretty=format:%B")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("erro ao obter mensagem do commit: %w", err)
	}
	return string(output), nil
}

// GetStagedFiles obtém os arquivos staged
func GetStagedFiles() ([]string, error) {
	cmd := exec.Command("git", "diff", "--cached", "--name-only")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("erro ao obter arquivos staged: %w", err)
	}

	if len(output) == 0 {
		return []string{}, nil
	}

	files := strings.Split(strings.TrimSpace(string(output)), "\n")
	return files, nil
}

// IsGitRepository verifica se estamos em um repositório git
func IsGitRepository() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	return cmd.Run() == nil
}

// FormatCommitMessage formata uma mensagem de commit
func FormatCommitMessage(commitType, scope, description string) string {
	if scope != "" {
		return fmt.Sprintf("%s(%s): %s", commitType, scope, description)
	}
	return fmt.Sprintf("%s: %s", commitType, description)
}

// GetCommitTypes retorna a lista de tipos válidos
func GetCommitTypes() []string {
	types := make([]string, len(ValidTypes))
	for i, t := range ValidTypes {
		types[i] = string(t)
	}
	return types
}
