package presenters

import "lanchonete/internal/domain/entities"

type ClienteDTO struct {
    Nome  string `json:"nome"`
    Email string `json:"email"`
    CPF   string `json:"cpf"` 
}

func formatCPF(cpf string) string {
    if len(cpf) != 11 {
        return cpf
    }
    return cpf[:3] + "." + cpf[3:6] + "." + cpf[6:9] + "-" + cpf[9:]
}

func NewClienteDTO(cliente *entities.Cliente) (*ClienteDTO) {
    if cliente == nil {
        return nil
    }
    
    return &ClienteDTO{
        Nome:  cliente.Nome,
        Email: cliente.Email,
		CPF:   formatCPF(cliente.CPF),
    }
}