package chainops

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type Status struct {
	Balances           map[string]uint
	TransactionMemPool []Transaction
	Dbfile             *os.File
}

type Transaction struct {
	De    string `json:"de"`
	Para  string `json:"para"`
	Valor uint   `json:"valor"`
	Data  string `json:"data"`
}

type Genesis struct {
	Balances map[string]uint `json:"balances"`
	Symbol   string          `json:"symbol"`

	ForkTIP1 uint64 `json:"fork_tip_1"`
}

// Aplicar a transação à cadeia
func (status *Status) apply(transaction Transaction) error {

	//Por ser uma recompença, nesse caso apenas adicionar mais credito a pessoa, sem tirar de outro usuario
	if transaction.Data == "reward" {
		status.Balances[transaction.De] += transaction.Valor
		return nil
	}

	if transaction.Valor > status.Balances[transaction.Para] {
		return fmt.Errorf("usuario nao tem creditos suficientes para realizar esta operação")
	}

	status.Balances[transaction.Para] -= transaction.Valor
	status.Balances[transaction.De] += transaction.Valor
	return nil
}

// Adicionar a transação ao pool de memoria da estrutura Status
func (status *Status) Adicionar(transaction Transaction) {
	err := status.apply(transaction)
	Check(err)
	status.TransactionMemPool = append(status.TransactionMemPool, transaction)
}

// Checando por erros no codigo e loga o erro, usar ao inves de iferr
func Check(err error) {
	if err != nil {
		log.Println(err)
	}
}

// Construindo o status atual a partir do genesis e de todos os balanços
func NewStatusFromDB() (*Status, error) {
	cwd, err := os.Getwd()
	Check(err)
	genFilePath := filepath.Join(cwd, "databases/json", "genesis.json")
	genesis := loadGenesis(genFilePath)

	balances := make(map[string]uint)
	for conta, balance := range genesis.Balances {
		balances[conta] = balance
	}
	transactionDBFilePath := filepath.Join(cwd, "databases", "transaction.db")
	transactionDB, err := os.OpenFile(transactionDBFilePath, os.O_APPEND|os.O_RDWR, 0600)
	Check(err)

	scanner := bufio.NewScanner(transactionDB)
	status := &Status{
		balances,
		make([]Transaction, 0),
		transactionDB,
	}

	//Alterando o balance das contas por cada linha q tem no banco de dados
	for scanner.Scan() {
		err := scanner.Err()
		Check(err)

		//Pegando o JSON e colocando pra dentro da estrutura
		var transaction Transaction
		json.Unmarshal(scanner.Bytes(), &transaction)
		err = status.apply(transaction)
		Check(err)
	}
	return status, nil
}

// Gravando as coisas no arquivo em disco
func (status *Status) Persistir() {

	//copia temporaria
	memPool := make([]Transaction, len(status.TransactionMemPool))
	copy(memPool, status.TransactionMemPool)

	for i := 0; i < len(memPool); i++ {
		transactionJson, err := json.Marshal(memPool[i])
		Check(err)

		_, err = status.Dbfile.Write(append(transactionJson, '\n'))
		Check(err)

		//Depois de gravar em disco, tira ele do Pool de Memória
		status.TransactionMemPool = status.TransactionMemPool[1:]
	}
}

func loadGenesis(path string) Genesis {
	content, err := os.ReadFile(path)
	Check(err)

	var loadedGenesis Genesis
	err = json.Unmarshal(content, &loadedGenesis)
	Check(err)

	return loadedGenesis
}
