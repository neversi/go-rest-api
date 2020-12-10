package database

import "gitlab.com/quybit/gexabyte/gexabyte_internship/go_abrd/models"

// RoleRepository ...
type RoleRepository struct {
	db *DataBase
}

// NewRoleRepository ...
func NewRoleRepository(db *DataBase) *RoleRepository {
	return &RoleRepository{
		db: db,
	}
}

// Create ...
func (rp *RoleRepository) Create(r *models.Role) error {
	currentDB := rp.db.Pdb

	if err := r.Validate(); err != nil {
		return err
	}
	
	currentDB.Create(r)

	return nil
}

func (rp *RoleRepository) Read(r *models.Role) (*models.Role, error) {
	currentDB := rp.db.Pdb
	newR := new(models.Role)
	res := currentDB.Model(&models.Role{}).Where("user_id = ?", r.UserID).First(&newR)

	if res.Error != nil {
		return nil, res.Error
	}
	return newR, nil
}

// Update role of the user
func (rp *RoleRepository) Update(r *models.Role) error {
	currentDB := rp.db.Pdb
	updatedR := new(models.Role)
	res := currentDB.Model(&models.Role{}).Where("user_id = ?", r.UserID).First(&updatedR)
	if res.Error != nil {
		return res.Error
	}
	
	updatedR.Role = r.Role

	res = currentDB.Model(&models.Role{}).Where("id = ?", r.ID).Save(&updatedR)

	if res.Error != nil {
		return res.Error
	}
	
	return nil
}

// Delete ...
func (rp *RoleRepository) Delete(r *models.Role) error {
	currentDB := rp.db.Pdb

	res := currentDB.Model(&models.Role{}).Where("id = ?", r.ID).Delete(&models.Role{})
	
	if res.Error != nil {
		return res.Error
	}
	
	return nil
}