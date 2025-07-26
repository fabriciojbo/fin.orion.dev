package commands

import (
	"fmt"
	"os"

	"fin.orion.dev/internal/commitlint"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// Comando para validar commits
var commitlintCmd = &cobra.Command{
	Use:   "commitlint [message]",
	Short: "Validar mensagem de commit",
	Long:  `Valida se uma mensagem de commit segue o padrão Conventional Commits.`,
	Args:  cobra.MaximumNArgs(1),
	RunE:  runCommitlint,
}

// Comando para validar último commit
var commitlintLastCmd = &cobra.Command{
	Use:   "commitlint-last",
	Short: "Validar último commit",
	Long:  `Valida se o último commit segue o padrão Conventional Commits.`,
	RunE:  runCommitlintLast,
}

// Comando para mostrar tipos válidos
var commitlintTypesCmd = &cobra.Command{
	Use:   "commitlint-types",
	Short: "Mostrar tipos válidos de commit",
	Long:  `Mostra todos os tipos válidos de commit para Conventional Commits.`,
	RunE:  runCommitlintTypes,
}

// Comando para formatar mensagem
var commitlintFormatCmd = &cobra.Command{
	Use:   "commitlint-format [type] [scope] [description]",
	Short: "Formatar mensagem de commit",
	Long:  `Formata uma mensagem de commit seguindo o padrão Conventional Commits.`,
	Args:  cobra.ExactArgs(3),
	RunE:  runCommitlintFormat,
}

// Comando para hook de pre-commit
var commitlintHookCmd = &cobra.Command{
	Use:   "commitlint-hook",
	Short: "Hook de pre-commit",
	Long:  `Hook para validar commits automaticamente (usado em .git/hooks/commit-msg).`,
	RunE:  runCommitlintHook,
}

func runCommitlint(cmd *cobra.Command, args []string) error {
	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)
	red := color.New(color.FgRed)
	yellow := color.New(color.FgYellow)

	var message string
	if len(args) > 0 {
		message = args[0]
	} else {
		// Ler da entrada padrão
		fmt.Print("Digite a mensagem de commit: ")
		fmt.Scanln(&message)
	}

	if message == "" {
		return fmt.Errorf("mensagem de commit é obrigatória")
	}

	_, _ = blue.Println("🔍 Validando mensagem de commit...")
	fmt.Println()

	// Criar validador
	validator := commitlint.NewValidator(nil)
	commit := validator.ValidateCommit(message)

	// Mostrar resultado
	if commit.IsValid {
		_, _ = green.Println("✅ Mensagem de commit válida!")
		fmt.Println()

		_, _ = blue.Println("📋 Detalhes:")
		fmt.Printf("  Tipo: %s\n", commit.Type)
		if commit.Scope != "" {
			fmt.Printf("  Escopo: %s\n", commit.Scope)
		}
		fmt.Printf("  Descrição: %s\n", commit.Description)
		if commit.Breaking {
			_, _ = yellow.Println("  ⚠️  Breaking Change detectado!")
		}
		if commit.Body != "" {
			fmt.Printf("  Body: %s\n", commit.Body)
		}
		if commit.Footer != "" {
			fmt.Printf("  Footer: %s\n", commit.Footer)
		}
	} else {
		_, _ = red.Println("❌ Mensagem de commit inválida!")
		fmt.Println()

		_, _ = red.Println("🚨 Erros encontrados:")
		for _, err := range commit.Errors {
			fmt.Printf("  • %s\n", err)
		}
		fmt.Println()

		_, _ = yellow.Println("💡 Exemplos de mensagens válidas:")
		fmt.Println("  feat: adicionar nova funcionalidade")
		fmt.Println("  fix(auth): corrigir validação de senha")
		fmt.Println("  docs: atualizar README")
		fmt.Println("  test(api): adicionar testes para endpoint")
		fmt.Println("  chore(deps): atualizar dependências")

		return fmt.Errorf("mensagem de commit inválida")
	}

	return nil
}

func runCommitlintLast(cmd *cobra.Command, args []string) error {
	blue := color.New(color.FgBlue)
	red := color.New(color.FgRed)

	// Verificar se estamos em um repositório git
	if !commitlint.IsGitRepository() {
		return fmt.Errorf("não estamos em um repositório git")
	}

	_, _ = blue.Println("🔍 Validando último commit...")
	fmt.Println()

	// Obter mensagem do último commit
	message, err := commitlint.GetCommitMessage()
	if err != nil {
		return fmt.Errorf("erro ao obter mensagem do último commit: %w", err)
	}

	if message == "" {
		return fmt.Errorf("não foi possível obter a mensagem do último commit")
	}

	// Validar mensagem
	validator := commitlint.NewValidator(nil)
	commit := validator.ValidateCommit(message)

	if commit.IsValid {
		_, _ = color.New(color.FgGreen).Println("✅ Último commit é válido!")
	} else {
		_, _ = red.Println("❌ Último commit é inválido!")
		fmt.Println()

		_, _ = red.Println("🚨 Erros encontrados:")
		for _, err := range commit.Errors {
			fmt.Printf("  • %s\n", err)
		}

		return fmt.Errorf("último commit é inválido")
	}

	return nil
}

func runCommitlintTypes(cmd *cobra.Command, args []string) error {
	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)

	_, _ = blue.Println("📋 Tipos válidos de commit:")
	fmt.Println()

	types := commitlint.GetCommitTypes()
	for _, t := range types {
		_, _ = green.Printf("  • %s\n", t)
	}

	fmt.Println()
	_, _ = blue.Println("💡 Descrições:")
	fmt.Println("  feat     - Nova funcionalidade")
	fmt.Println("  fix      - Correção de bug")
	fmt.Println("  docs     - Documentação")
	fmt.Println("  style    - Formatação de código")
	fmt.Println("  refactor - Refatoração de código")
	fmt.Println("  test     - Testes")
	fmt.Println("  chore    - Manutenção")
	fmt.Println("  perf     - Melhorias de performance")
	fmt.Println("  ci       - Integração contínua")
	fmt.Println("  build    - Build do sistema")
	fmt.Println("  revert   - Reverter commit")

	return nil
}

func runCommitlintFormat(cmd *cobra.Command, args []string) error {
	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)

	commitType := args[0]
	scope := args[1]
	description := args[2]

	_, _ = blue.Println("🔧 Formatando mensagem de commit...")
	fmt.Println()

	formatted := commitlint.FormatCommitMessage(commitType, scope, description)
	_, _ = green.Printf("✅ Mensagem formatada: %s\n", formatted)

	return nil
}

func runCommitlintHook(cmd *cobra.Command, args []string) error {
	red := color.New(color.FgRed)

	// Verificar se estamos em um repositório git
	if !commitlint.IsGitRepository() {
		return fmt.Errorf("não estamos em um repositório git")
	}

	// Obter arquivo de mensagem de commit
	commitMsgFile := os.Getenv("COMMIT_EDITMSG")
	if commitMsgFile == "" {
		commitMsgFile = ".git/COMMIT_EDITMSG"
	}

	// Ler mensagem do arquivo
	content, err := os.ReadFile(commitMsgFile)
	if err != nil {
		return fmt.Errorf("erro ao ler arquivo de mensagem: %w", err)
	}

	message := string(content)

	// Validar mensagem
	validator := commitlint.NewValidator(nil)
	commit := validator.ValidateCommit(message)

	if !commit.IsValid {
		_, _ = red.Println("❌ Commit rejeitado!")
		fmt.Println()

		_, _ = red.Println("🚨 Erros encontrados:")
		for _, err := range commit.Errors {
			fmt.Printf("  • %s\n", err)
		}
		fmt.Println()

		_, _ = red.Println("💡 Use o formato: type(scope): description")
		_, _ = red.Println("   Exemplo: feat(auth): adicionar autenticação JWT")

		return fmt.Errorf("mensagem de commit inválida")
	}

	return nil
}
