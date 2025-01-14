//go:build linux && cgo && !agent
// +build linux,cgo,!agent

package db

// The code below was generated by lxd-generate - DO NOT EDIT!

import (
	"fmt"

	"github.com/lxc/lxd/lxd/db/cluster"
	"github.com/lxc/lxd/lxd/db/query"
	"github.com/lxc/lxd/shared/api"
)

var _ = api.ServerEnvironment{}

var certificateProjectObjects = cluster.RegisterStmt(`
SELECT certificates_projects.certificate_id, certificates_projects.project_id
  FROM certificates_projects
  ORDER BY certificates_projects.certificate_id
`)

var certificateProjectCreate = cluster.RegisterStmt(`
INSERT INTO certificates_projects (certificate_id, project_id)
  VALUES (?, ?)
`)

var certificateProjectDeleteByCertificateID = cluster.RegisterStmt(`
DELETE FROM certificates_projects WHERE certificate_id = ?
`)

// GetCertificateProjects returns all available certificate_projects.
// generator: certificate_project GetMany
func (c *ClusterTx) GetCertificateProjects() (map[int][]int, error) {
	var err error

	// Result slice.
	objects := make([]CertificateProject, 0)

	stmt := c.stmt(certificateProjectObjects)
	args := []interface{}{}

	// Dest function for scanning a row.
	dest := func(i int) []interface{} {
		objects = append(objects, CertificateProject{})
		return []interface{}{
			&objects[i].CertificateID,
			&objects[i].ProjectID,
		}
	}

	// Select.
	err = query.SelectObjects(stmt, dest, args...)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch from \"certificates_projects\" table: %w", err)
	}

	resultMap := map[int][]int{}
	for _, object := range objects {
		resultMap[object.CertificateID] = append(resultMap[object.CertificateID], object.ProjectID)
	}

	return resultMap, nil
}

// DeleteCertificateProjects deletes the certificate_project matching the given key parameters.
// generator: certificate_project DeleteMany
func (c *ClusterTx) DeleteCertificateProjects(object Certificate) error {
	stmt := c.stmt(certificateProjectDeleteByCertificateID)
	result, err := stmt.Exec(int(object.ID))
	if err != nil {
		return fmt.Errorf("Delete \"certificates_projects\" entry failed: %w", err)
	}

	_, err = result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Fetch affected rows: %w", err)
	}

	return nil
}

// CreateCertificateProject adds a new certificate_project to the database.
// generator: certificate_project Create
func (c *ClusterTx) CreateCertificateProject(object CertificateProject) (int64, error) {
	args := make([]interface{}, 2)

	// Populate the statement arguments.
	args[0] = object.CertificateID
	args[1] = object.ProjectID

	// Prepared statement to use.
	stmt := c.stmt(certificateProjectCreate)

	// Execute the statement.
	result, err := stmt.Exec(args...)
	if err != nil {
		return -1, fmt.Errorf("Failed to create \"certificates_projects\" entry: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("Failed to fetch \"certificates_projects\" entry ID: %w", err)
	}

	return id, nil
}

// UpdateCertificateProjects updates the certificate_project matching the given key parameters.
// generator: certificate_project Update
func (c *ClusterTx) UpdateCertificateProjects(object Certificate) error {
	// Delete current entry.
	err := c.DeleteCertificateProjects(object)
	if err != nil {
		return err
	}

	// Insert new entries.
	for _, key := range object.Projects {
		refID, err := c.GetProjectID(key)
		if err != nil {
			return err
		}

		certificateProject := CertificateProject{CertificateID: object.ID, ProjectID: int(refID)}
		_, err = c.CreateCertificateProject(certificateProject)
		if err != nil {
			return err
		}

		return nil
	}
	return nil
}
