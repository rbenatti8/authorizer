package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/authorizer/internal/core/service"
	"github.com/authorizer/internal/driven/database"
	"github.com/authorizer/internal/driven/repository"
	"github.com/stretchr/testify/assert"
)

//func TestNewHandler(t *testing.T) {
//	_ = json.NewDecoder(strings.NewReader(""))
//	expect := Handler{}
//	result := NewHandler()
//	assert.Equal(t, expect, result)
//}

func TestHandler_Handle(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "criando uma conta com sucesso",
			input:    "{\"account\":{\"active-card\":false,\"available-limit\":750}}\n",
			expected: "{\"account\":{\"active-card\":false,\"available-limit\":750},\"violations\":[]}\n",
		},
		{
			name:     "Criando uma conta que viola a lógica do Autorizador",
			input:    "{\"account\":{\"active-card\":false,\"available-limit\":175}}\n{\"account\":{\"active-card\":false,\"available-limit\":350}}\n",
			expected: "{\"account\":{\"active-card\":false,\"available-limit\":175},\"violations\":[]}\n{\"account\":{\"active-card\":false,\"available-limit\":175},\"violations\":[\"account-already-initialized\"]}\n",
		},
		{
			name:     "Processando uma transação com sucesso",
			input:    "{\"account\":{\"active-card\":true,\"available-limit\":100}}\n{\"transaction\":{\"merchant\":\"Burger King\",\"amount\":20,\"time\":\"2019-02-13T11:00:00.000Z\"}}",
			expected: "{\"account\":{\"active-card\":true,\"available-limit\":100},\"violations\":[]}\n{\"account\":{\"active-card\":true,\"available-limit\":80},\"violations\":[]}\n",
		},
		{
			name:     "Processando uma transação que viola a lógica account-not-initialized",
			input:    "{\"transaction\":{\"merchant\":\"Uber Eats\",\"amount\":25,\"time\":\"2020-12-01T11:07:00.000Z\"}}\n{\"account\":{\"active-card\":true,\"available-limit\":225}}\n{\"transaction\":{\"merchant\":\"Uber Eats\",\"amount\":25,\"time\":\"2020-12-01T11:07:00.000Z\"}}\n",
			expected: "{\"account\":{},\"violations\":[\"account-not-initialized\"]}\n{\"account\":{\"active-card\":true,\"available-limit\":225},\"violations\":[]}\n{\"account\":{\"active-card\":true,\"available-limit\":200},\"violations\":[]}\n",
		},
		{
			name:     "Processando uma transação que viola a lógica card-not-active",
			input:    "{\"account\":{\"active-card\":false,\"available-limit\":100}}\n{\"transaction\":{\"merchant\":\"Burger King\",\"amount\":20,\"time\":\"2019-02-13T11:00:00.000Z\"}}\n{\"transaction\":{\"merchant\":\"Habbib's\",\"amount\":15,\"time\":\"2019-02-13T11:15:00.000Z\"}}\n",
			expected: "{\"account\":{\"active-card\":false,\"available-limit\":100},\"violations\":[]}\n{\"account\":{\"active-card\":false,\"available-limit\":100},\"violations\":[\"card-not-active\"]}\n{\"account\":{\"active-card\":false,\"available-limit\":100},\"violations\":[\"card-not-active\"]}\n",
		},
		{
			name:     "Processando uma transação que viola a lógica insufficient-limit",
			input:    "{\"account\":{\"active-card\":true,\"available-limit\":1000}}\n{\"transaction\":{\"merchant\":\"Vivara\",\"amount\":1250,\"time\":\"2019-02-13T11:00:00.000Z\"}}\n{\"transaction\":{\"merchant\":\"Samsung\",\"amount\":2500,\"time\":\"2019-02-13T11:00:01.000Z\"}}\n{\"transaction\":{\"merchant\":\"Nike\",\"amount\":800,\"time\":\"2019-02-13T11:01:01.000Z\"}}\n",
			expected: "{\"account\":{\"active-card\":true,\"available-limit\":1000},\"violations\":[]}\n{\"account\":{\"active-card\":true,\"available-limit\":1000},\"violations\":[\"insufficient-limit\"]}\n{\"account\":{\"active-card\":true,\"available-limit\":1000},\"violations\":[\"insufficient-limit\"]}\n{\"account\":{\"active-card\":true,\"available-limit\":200},\"violations\":[]}\n",
		},
		{
			name:     "Processando uma transação que viola a lógica high-frequency-small-interval",
			input:    "{\"account\":{\"active-card\":true,\"available-limit\":100}}\n{\"transaction\":{\"merchant\":\"Burger King\",\"amount\":20,\"time\":\"2019-02-13T11:00:00.000Z\"}}\n{\"transaction\":{\"merchant\":\"Habbib's\",\"amount\":20,\"time\":\"2019-02-13T11:00:01.000Z\"}}\n{\"transaction\":{\"merchant\":\"McDonald's\",\"amount\":20,\"time\":\"2019-02-13T11:01:01.000Z\"}}\n{\"transaction\":{\"merchant\":\"Subway\",\"amount\":20,\"time\":\"2019-02-13T11:01:31.000Z\"}}\n{\"transaction\":{\"merchant\":\"Burger King\",\"amount\":10,\"time\":\"2019-02-13T12:00:00.000Z\"}}\n",
			expected: "{\"account\":{\"active-card\":true,\"available-limit\":100},\"violations\":[]}\n{\"account\":{\"active-card\":true,\"available-limit\":80},\"violations\":[]}\n{\"account\":{\"active-card\":true,\"available-limit\":60},\"violations\":[]}\n{\"account\":{\"active-card\":true,\"available-limit\":40},\"violations\":[]}\n{\"account\":{\"active-card\":true,\"available-limit\":40},\"violations\":[\"high-frequency-small-interval\"]}\n{\"account\":{\"active-card\":true,\"available-limit\":30},\"violations\":[]}\n",
		},
		{
			name:     "Processando uma transação que viola a lógica doubled-transaction",
			input:    "{\"account\":{\"active-card\":true,\"available-limit\":100}}\n{\"transaction\":{\"merchant\":\"Burger King\",\"amount\":20,\"time\":\"2019-02-13T11:00:00.000Z\"}}\n{\"transaction\":{\"merchant\":\"McDonald's\",\"amount\":10,\"time\":\"2019-02-13T11:00:01.000Z\"}}\n{\"transaction\":{\"merchant\":\"Burger King\",\"amount\":20,\"time\":\"2019-02-13T11:00:02.000Z\"}}\n{\"transaction\":{\"merchant\":\"Burger King\",\"amount\":15,\"time\":\"2019-02-13T11:00:03.000Z\"}}\n",
			expected: "{\"account\":{\"active-card\":true,\"available-limit\":100},\"violations\":[]}\n{\"account\":{\"active-card\":true,\"available-limit\":80},\"violations\":[]}\n{\"account\":{\"active-card\":true,\"available-limit\":70},\"violations\":[]}\n{\"account\":{\"active-card\":true,\"available-limit\":70},\"violations\":[\"doubled-transaction\"]}\n{\"account\":{\"active-card\":true,\"available-limit\":55},\"violations\":[]}\n",
		},
		{
			name:     "Processando transações que violam multiplas lógicas",
			input:    "{\"account\":{\"active-card\":true,\"available-limit\":100}}\n{\"transaction\":{\"merchant\":\"McDonald's\",\"amount\":10,\"time\":\"2019-02-13T11:00:01.000Z\"}}\n{\"transaction\":{\"merchant\":\"Burger King\",\"amount\":20,\"time\":\"2019-02-13T11:00:02.000Z\"}}\n{\"transaction\":{\"merchant\":\"Burger King\",\"amount\":5,\"time\":\"2019-02-13T11:00:07.000Z\"}}\n{\"transaction\":{\"merchant\":\"Burger King\",\"amount\":5,\"time\":\"2019-02-13T11:00:08.000Z\"}}\n{\"transaction\":{\"merchant\":\"Burger King\",\"amount\":150,\"time\":\"2019-02-13T11:00:18.000Z\"}}\n{\"transaction\":{\"merchant\":\"Burger King\",\"amount\":190,\"time\":\"2019-02-13T11:00:22.000Z\"}}\n{\"transaction\":{\"merchant\":\"Burger King\",\"amount\":15,\"time\":\"2019-02-13T12:00:27.000Z\"}}\n",
			expected: "{\"account\":{\"active-card\":true,\"available-limit\":100},\"violations\":[]}\n{\"account\":{\"active-card\":true,\"available-limit\":90},\"violations\":[]}\n{\"account\":{\"active-card\":true,\"available-limit\":70},\"violations\":[]}\n{\"account\":{\"active-card\":true,\"available-limit\":65},\"violations\":[]}\n{\"account\":{\"active-card\":true,\"available-limit\":65},\"violations\":[\"high-frequency-small-interval\",\"doubled-transaction\"]}\n{\"account\":{\"active-card\":true,\"available-limit\":65},\"violations\":[\"insufficient-limit\",\"high-frequency-small-interval\"]}\n{\"account\":{\"active-card\":true,\"available-limit\":65},\"violations\":[\"insufficient-limit\",\"high-frequency-small-interval\"]}\n{\"account\":{\"active-card\":true,\"available-limit\":50},\"violations\":[]}\n",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			var stdout bytes.Buffer
			stdin := strings.NewReader(tt.input)

			db := database.NewInMemoryDB()

			accountRepo := repository.NewAccountRepository(db)

			as := service.NewAccount(accountRepo)
			ts := service.NewTransaction(accountRepo)

			handler := NewHandler(as, ts)

			_ = handler.Handle(stdin, &stdout)
			assert.Equal(t, tt.expected, stdout.String())
		})
	}
}
