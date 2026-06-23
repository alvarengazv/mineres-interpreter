package utils

import (
	"os"
	"path/filepath"
	"testing"
)

// expectPanic executa fn e verifica se houve um panic do tipo CompilerError.
// Retorna a mensagem do erro capturado, ou faz t.Fatal se nenhum panic ocorreu.
func expectPanic(t *testing.T, fn func()) string {
	t.Helper()
	old := ExitOnError
	ExitOnError = false
	defer func() { ExitOnError = old }()

	var msg string
	func() {
		defer func() {
			r := recover()
			if r == nil {
				t.Fatal("expected panic but none occurred")
			}
			ce, ok := r.(CompilerError)
			if !ok {
				t.Fatalf("expected CompilerError, got %T: %v", r, r)
			}
			msg = ce.Message
		}()
		fn()
	}()
	return msg
}

// ---------- ArquivoExiste ----------

func TestArquivoExiste_Existe(t *testing.T) {
	// Cria arquivo temporário no diretório do projeto
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.uai")
	if err := os.WriteFile(tmpFile, []byte("conteudo"), 0644); err != nil {
		t.Fatal(err)
	}
	if !ArquivoExiste(tmpFile) {
		t.Errorf("ArquivoExiste(%q) = false, want true", tmpFile)
	}
}

func TestArquivoExiste_NaoExiste(t *testing.T) {
	if ArquivoExiste("/caminho/que/nao/existe/arquivo.uai") {
		t.Error("ArquivoExiste should return false for non-existent file")
	}
}

// ---------- LerArquivo ----------

func TestLerArquivo_Valido(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "hello.uai")
	conteudo := "bora_cumpade main() simbora cabo"
	if err := os.WriteFile(tmpFile, []byte(conteudo), 0644); err != nil {
		t.Fatal(err)
	}
	resultado := LerArquivo(tmpFile)
	if resultado != conteudo {
		t.Errorf("LerArquivo = %q, want %q", resultado, conteudo)
	}
}

func TestLerArquivo_ArquivoInexistente(t *testing.T) {
	msg := expectPanic(t, func() {
		LerArquivo("/arquivo/inexistente.uai")
	})
	if msg == "" {
		t.Error("expected error message, got empty string")
	}
}
