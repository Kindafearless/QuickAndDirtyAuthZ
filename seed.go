package main

func (a *App) SeedData() error {

	err := a.DB.WriteToDatabase(a.DataTableName, "jbob", "joe@bobsite.com")
	if err != nil { return err }
	err = a.DB.WriteToDatabase(a.DataTableName, "bbob", "billybob2304@gmail.com")
	if err != nil { return err }
	err = a.DB.WriteToDatabase(a.DataTableName, "bjoe", "bjoe@aol.com")
	if err != nil { return err }

	return nil
}

func (a *App) SeedUsers() error {

	err := a.DB.WriteToDatabase(a.UserTableName, "jbob", "admin")
	if err != nil { return err }
	err = a.DB.WriteToDatabase(a.UserTableName, "bbob", "developer")
	if err != nil { return err }
	err = a.DB.WriteToDatabase(a.UserTableName, "bjoe", "developer")
	if err != nil { return err }

	return nil
}