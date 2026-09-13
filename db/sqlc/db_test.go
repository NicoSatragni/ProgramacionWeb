package db

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
	"testing"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}
	user := os.Getenv("DB_USER")
	if user == "" {
		user = "nSatragni"
	}
	pass := os.Getenv("DB_PASSWORD")
	if pass == "" {
		pass = "12345"
	}
	nom := os.Getenv("DB_NAME")
	if nom == "" {
		nom = "apirest"
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, nom)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("error conectando a la bd: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("No hay conexión con la bd: %v", err)
	}
	fmt.Println("Conexión a la base de datos exitosa.")

	testQueries = New(db)

	code := m.Run()

	db.Close()
	os.Exit(code)
}

// Tests
//

// Probamos la inserción y recuperación de usuarios.
func TestInsertGetUser(t *testing.T) {
	ctx := context.Background()

	name := "UsuarioTest"
	newUser, err := testQueries.CreateUser(ctx, name)
	if err != nil {
		t.Fatalf("Error: No se pudo crear el usuario. No tiene sentido seguir haciendo pruebas: %v", err)
	}

	if newUser.ID == 0 {
		t.Errorf("El usuario deberia tener ID mayor a 0, se obtuvo: %d", newUser.ID)
	}

	if newUser.Name != name {
		t.Errorf("El nombre del usuario no coincide: '%s', y se esperaba: '%s'", newUser.Name, name)
	}

	if !newUser.CreatedAt.Valid {
		t.Errorf("La fecha y/o hora no son válidos: %v", newUser.CreatedAt)
	}

	t.Logf("Usuario creado con éxito: ID=%d, Nombre=%s", newUser.ID, newUser.Name)

	// Probamos recuperar usuario
	getUser, err := testQueries.GetUser(ctx, newUser.ID)
	if err != nil {
		t.Fatalf("No se recuperó el usuario correctamente: %v", err)
	}

	if getUser.Name != newUser.Name {
		t.Errorf("No se recuperó el usuario correctamente: ID=%d, Nombre=%s. Esperado: ID=%d, Nombre=%s",
			getUser.ID, getUser.Name, newUser.ID, newUser.Name)
	}
	t.Logf("Se recuperó el usuario correctamente: ID=%d", getUser.ID)

	// Probamos usuario aleatorio
	randomUs, err := testQueries.GetRandomUser(ctx)
	if err != nil {
		t.Errorf("No se pudo obtener un usuario de forma aleatoria: %v", err)
	}
	if randomUs.ID <= 0 || randomUs.Name == "" || !randomUs.CreatedAt.Valid {
		t.Errorf("El usuario que se recuperó no está completo: ID=%d, Nombre=%s, CreatedAt=%v",
			randomUs.ID, randomUs.Name, randomUs.CreatedAt)
	}
	t.Logf("Usuario aleatorio obtenido: ID=%d, Nombre=%s", randomUs.ID, randomUs.Name)
}

func TestListUsers(t *testing.T) {
	ctx := context.Background()

	// Creamos un usuario para garantizar que la tabla no esté vacía en este test
	_, err := testQueries.CreateUser(ctx, "UsuarioListado")
	if err != nil {
		t.Fatalf("Error preparando usuario para TestListUsers: %v", err)
	}

	// Probamos el contador de usuarios
	cntUsers, err := testQueries.CountUsers(ctx)
	if err != nil {
		t.Errorf("No se contaron correctamente los usuarios: %v", err)
	}
	if cntUsers <= 0 {
		t.Errorf("Se esperaba al menos un usuario, se obtuvo: %d", cntUsers)
	}

	// Obtener todos los usuarios
	listUsers, err := testQueries.ListUsers(ctx)
	if err != nil {
		t.Fatalf("No se pudieron listar los usuarios: %v", err)
	}

	if len(listUsers) == 0 {
		t.Errorf("Se esperaba al menos 1 usuario en la lista")
	}

	t.Logf("Total de usuarios listados: %d", len(listUsers))
}

// Probamos el CRUD y ciclo de vida de los artículos de compra
func TestShoppingItemLifecycle(t *testing.T) {
	ctx := context.Background()

	// 1. Crear un usuario comprador
	payer, err := testQueries.CreateUser(ctx, "Comprador")
	if err != nil {
		t.Fatalf("Error creando usuario comprador: %v", err)
	}

	// 2. Crear un artículo pendiente
	item, err := testQueries.CreateShoppingItem(ctx, CreateShoppingItemParams{
		Title:    "Leche y Café",
		Quantity: 2,
	})
	if err != nil {
		t.Fatalf("Error creando artículo de compra: %v", err)
	}

	if item.ID <= 0 || item.Title != "Leche y Café" || item.Quantity != 2 || item.IsPurchased {
		t.Fatalf("El artículo no se creó con los valores esperados: %+v", item)
	}

	// 3. Actualizar datos del artículo (UpdateShoppingItem)
	err = testQueries.UpdateShoppingItem(ctx, UpdateShoppingItemParams{
		ID:       item.ID,
		Title:    "Café de especialidad",
		Quantity: 3,
	})
	if err != nil {
		t.Fatalf("Error actualizando artículo: %v", err)
	}

	updatedItem, err := testQueries.GetShoppingItem(ctx, item.ID)
	if err != nil {
		t.Fatalf("Error recuperando artículo actualizado: %v", err)
	}
	if updatedItem.Title != "Café de especialidad" || updatedItem.Quantity != 3 {
		t.Errorf("La actualización no se aplicó correctamente: %+v", updatedItem)
	}

	// 4. Marcar como comprado (MarkItemAsPurchased)
	purchasedItem, err := testQueries.MarkItemAsPurchased(ctx, MarkItemAsPurchasedParams{
		ID:           item.ID,
		Price:        "1500.50",
		PaidByUserID: sql.NullInt32{Int32: payer.ID, Valid: true},
	})
	if err != nil {
		t.Fatalf("Error marcando artículo como comprado: %v", err)
	}
	if !purchasedItem.IsPurchased || purchasedItem.Price != "1500.50" || purchasedItem.PaidByUserID.Int32 != payer.ID {
		t.Errorf("El artículo comprado no tiene los datos correctos: %+v", purchasedItem)
	}

	// 5. Listar pendientes y comprados
	pending, err := testQueries.ListPendingItems(ctx)
	if err != nil {
		t.Fatalf("Error al listar pendientes: %v", err)
	}
	for _, p := range pending {
		if p.ID == item.ID {
			t.Errorf("El artículo con ID %d sigue apareciendo como pendiente", item.ID)
		}
	}

	purchasedList, err := testQueries.ListPurchasedItems(ctx)
	if err != nil {
		t.Fatalf("Error al listar comprados: %v", err)
	}
	if len(purchasedList) == 0 {
		t.Errorf("Se esperaba al menos 1 artículo comprado")
	}

	// 6. Desmarcar como comprado (UnmarkItemAsPurchased)
	unmarkedItem, err := testQueries.UnmarkItemAsPurchased(ctx, item.ID)
	if err != nil {
		t.Fatalf("Error al desmarcar artículo: %v", err)
	}
	if unmarkedItem.IsPurchased || unmarkedItem.PaidByUserID.Valid {
		t.Errorf("El artículo no se desmarcó correctamente: %+v", unmarkedItem)
	}

	// 7. Borrar artículo (DeleteShoppingItem)
	err = testQueries.DeleteShoppingItem(ctx, item.ID)
	if err != nil {
		t.Fatalf("Error al borrar artículo: %v", err)
	}

	_, err = testQueries.GetShoppingItem(ctx, item.ID)
	if err == nil {
		t.Errorf("Se esperaba error al buscar un artículo eliminado")
	}
}

// Probamos la división de gastos entre usuarios (item_splits)
func TestItemSplits(t *testing.T) {
	ctx := context.Background()

	// 1. Crear usuarios y artículo
	u1, err := testQueries.CreateUser(ctx, "Inquilino 1")
	if err != nil {
		t.Fatalf("Error creando usuario: %v", err)
	}
	u2, err := testQueries.CreateUser(ctx, "Inquilino 2")
	if err != nil {
		t.Fatalf("Error creando usuario: %v", err)
	}

	item, err := testQueries.CreateShoppingItem(ctx, CreateShoppingItemParams{
		Title:    "Detergente",
		Quantity: 1,
	})
	if err != nil {
		t.Fatalf("Error creando artículo: %v", err)
	}

	// 2. Asignar división a u1 y u2
	err = testQueries.CreateItemSplit(ctx, CreateItemSplitParams{
		ItemID: item.ID,
		UserID: u1.ID,
	})
	if err != nil {
		t.Fatalf("Error creando split individual: %v", err)
	}

	err = testQueries.CreateItemSplit(ctx, CreateItemSplitParams{
		ItemID: item.ID,
		UserID: u2.ID,
	})
	if err != nil {
		t.Fatalf("Error creando split individual: %v", err)
	}

	// 3. Contar y listar splits por artículo
	count, err := testQueries.CountItemSplitsByItem(ctx, item.ID)
	if err != nil {
		t.Fatalf("Error contando splits: %v", err)
	}
	if count != 2 {
		t.Errorf("Se esperaban 2 splits, se obtuvieron: %d", count)
	}

	splits, err := testQueries.ListItemSplitsByItem(ctx, item.ID)
	if err != nil {
		t.Fatalf("Error listando splits: %v", err)
	}
	if len(splits) != 2 {
		t.Errorf("Se esperaban 2 elementos en la lista de splits, se obtuvieron: %d", len(splits))
	}

	// 4. Eliminar un split
	err = testQueries.DeleteItemSplit(ctx, DeleteItemSplitParams{
		ItemID: item.ID,
		UserID: u1.ID,
	})
	if err != nil {
		t.Fatalf("Error al eliminar un split: %v", err)
	}

	countAfterDelete, err := testQueries.CountItemSplitsByItem(ctx, item.ID)
	if err != nil {
		t.Fatalf("Error contando splits tras borrado: %v", err)
	}
	if countAfterDelete != 1 {
		t.Errorf("Se esperaba 1 split restante, se obtuvieron: %d", countAfterDelete)
	}

	// 5. Eliminar todos los splits del artículo
	err = testQueries.DeleteItemSplitsByItem(ctx, item.ID)
	if err != nil {
		t.Fatalf("Error borrando todos los splits: %v", err)
	}
	countFinal, _ := testQueries.CountItemSplitsByItem(ctx, item.ID)
	if countFinal != 0 {
		t.Errorf("Se esperaban 0 splits, se obtuvieron: %d", countFinal)
	}
}

// Probamos los acuerdos/transferencias entre usuarios (Settlements)
func TestSettlementsAndBalances(t *testing.T) {
	ctx := context.Background()

	// 1. Crear pagador y receptor
	payer, err := testQueries.CreateUser(ctx, "Pagador Deuda")
	if err != nil {
		t.Fatalf("Error creando usuario pagador: %v", err)
	}
	receiver, err := testQueries.CreateUser(ctx, "Receptor Deuda")
	if err != nil {
		t.Fatalf("Error creando usuario receptor: %v", err)
	}

	// 2. Crear un settlement (pago de deuda)
	amount := "2500.00"
	settlement, err := testQueries.CreateSettlement(ctx, CreateSettlementParams{
		PayerID:    payer.ID,
		ReceiverID: receiver.ID,
		Amount:     amount,
	})
	if err != nil {
		t.Fatalf("Error registrando pago de deuda: %v", err)
	}
	if settlement.ID <= 0 || settlement.Amount != amount {
		t.Errorf("El settlement no coincide con lo esperado: %+v", settlement)
	}

	// 3. Obtener el settlement con nombres (GetSettlement)
	fetched, err := testQueries.GetSettlement(ctx, settlement.ID)
	if err != nil {
		t.Fatalf("Error recuperando settlement: %v", err)
	}
	if fetched.PayerName != payer.Name || fetched.ReceiverName != receiver.Name {
		t.Errorf("Los nombres de los involucrados no coinciden: pagador=%s, receptor=%s",
			fetched.PayerName, fetched.ReceiverName)
	}

	// 4. Probar balances acumulados
	totalPaid, err := testQueries.GetTotalSettlementsPaidByUser(ctx, payer.ID)
	if err != nil {
		t.Fatalf("Error obteniendo total pagado en deudas: %v", err)
	}
	if totalPaid != amount {
		t.Errorf("Total pagado esperado %s, obtenido %s", amount, totalPaid)
	}

	totalReceived, err := testQueries.GetTotalSettlementsReceivedByUser(ctx, receiver.ID)
	if err != nil {
		t.Fatalf("Error obteniendo total recibido en deudas: %v", err)
	}
	if totalReceived != amount {
		t.Errorf("Total recibido esperado %s, obtenido %s", amount, totalReceived)
	}

	// 5. Borrar el settlement
	err = testQueries.DeleteSettlement(ctx, settlement.ID)
	if err != nil {
		t.Fatalf("Error eliminando settlement: %v", err)
	}
}
