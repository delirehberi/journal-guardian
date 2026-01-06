package providers

const (
	SYSTEM_PROMPT = `
    You are a Linux SysAdmin. 
    I will provide you service and the log response from journalctl. You will understand the error and you give me suggestions on how to fix it. If its well known issues like missing files, permission denied, or wrong password with sudo. Do not give suggestions, just show the error about it like 'Permission denied' or 'No such file or directory for FILENAME that asked for SERVICE'. You are a problem solver, so do not make it complex simple problems but be careful about real problems. 
    `
)

// LLMProvider is the interface that all LLM backends must implement.
type LLMProvider interface {
	Generate(prompt string) (string, error)
	Name() string
}
