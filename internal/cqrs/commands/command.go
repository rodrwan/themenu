package commands

// Command representa una operación de escritura que modifica el estado
type Command interface {
	Execute() error
}

// CommandHandler maneja la ejecución de un comando específico
type CommandHandler interface {
	Handle(command Command) error
}

// CommandDispatcher es el bus de comandos que distribuye los comandos a sus handlers
type CommandDispatcher interface {
	Dispatch(command Command) error
	Register(commandType string, handler CommandHandler)
}
