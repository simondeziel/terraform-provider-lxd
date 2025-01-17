package instance_test

import (
	"fmt"
	"testing"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/terraform-lxd/terraform-provider-lxd/internal/acctest"
)

func TestAccInstanceFile_basic(t *testing.T) {
	instanceName := petname.Generate(2, "-")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceFile_content(instanceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("lxd_instance.instance1", "name", instanceName),
					resource.TestCheckResourceAttr("lxd_instance.instance1", "status", "Running"),
					resource.TestCheckResourceAttr("lxd_instance_file.file1", "instance", instanceName),
					resource.TestCheckResourceAttr("lxd_instance_file.file1", "content", "Hello, World!\n"),
					resource.TestCheckNoResourceAttr("lxd_instance_file.file1", "source_path"),
					resource.TestCheckResourceAttr("lxd_instance_file.file1", "target_path", "/foo/bar.txt"),
					resource.TestCheckResourceAttr("lxd_instance_file.file1", "create_directories", "true"),
					resource.TestCheckResourceAttr("lxd_instance_file.file1", "resource_id", fmt.Sprintf(":%s:/foo/bar.txt", instanceName)),
				),
			},
			{
				// Ensure no changes happen.
				Config: testAccInstanceFile_content(instanceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("lxd_instance.instance1", "name", instanceName),
					resource.TestCheckResourceAttr("lxd_instance.instance1", "status", "Running"),
					resource.TestCheckResourceAttr("lxd_instance_file.file1", "resource_id", fmt.Sprintf(":%s:/foo/bar.txt", instanceName)),
				),
			},
			{
				// Upload file from source instead of content.
				// This should recreate the file.
				Config: testAccInstanceFile_sourcePath(instanceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("lxd_instance.instance1", "name", instanceName),
					resource.TestCheckResourceAttr("lxd_instance.instance1", "status", "Running"),
					resource.TestCheckResourceAttr("lxd_instance_file.file1", "instance", instanceName),
					resource.TestCheckNoResourceAttr("lxd_instance_file.file1", "content"),
					resource.TestCheckResourceAttr("lxd_instance_file.file1", "source_path", "../acctest/fixtures/test-file.txt"),
					resource.TestCheckResourceAttr("lxd_instance_file.file1", "target_path", "/foo/bar.txt"),
					resource.TestCheckResourceAttr("lxd_instance_file.file1", "create_directories", "true"),
					resource.TestCheckResourceAttr("lxd_instance_file.file1", "resource_id", fmt.Sprintf(":%s:/foo/bar.txt", instanceName)),
				),
			},
		},
	})
}

func TestAccInstanceFile_project(t *testing.T) {
	projectName := petname.Name()
	instanceName := petname.Generate(2, "-")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceFile_project(projectName, instanceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("lxd_project.project1", "name", projectName),
					resource.TestCheckResourceAttr("lxd_instance.instance1", "name", instanceName),
					resource.TestCheckResourceAttr("lxd_instance.instance1", "project", projectName),
					resource.TestCheckResourceAttr("lxd_instance.instance1", "status", "Running"),
					resource.TestCheckResourceAttr("lxd_instance_file.file1", "instance", instanceName),
					resource.TestCheckResourceAttr("lxd_instance_file.file1", "project", projectName),
					resource.TestCheckResourceAttr("lxd_instance_file.file1", "source_path", "../acctest/fixtures/test-file.txt"),
					resource.TestCheckResourceAttr("lxd_instance_file.file1", "target_path", "/foo/bar.txt"),
					resource.TestCheckResourceAttr("lxd_instance_file.file1", "create_directories", "true"),
					resource.TestCheckResourceAttr("lxd_instance_file.file1", "resource_id", fmt.Sprintf(":%s:/foo/bar.txt", instanceName)),
				),
			},
		},
	})
}

func testAccInstanceFile_content(name string) string {
	return fmt.Sprintf(`
resource "lxd_instance" "instance1" {
  name  = "%s"
  image = "%s"
}

resource "lxd_instance_file" "file1" {
  instance           = lxd_instance.instance1.name
  content            = "Hello, World!\n"
  target_path        = "/foo/bar.txt"
  create_directories = true
}
	`, name, acctest.TestImage)
}

func testAccInstanceFile_sourcePath(name string) string {
	return fmt.Sprintf(`
resource "lxd_instance" "instance1" {
  name  = "%s"
  image = "%s"
}

resource "lxd_instance_file" "file1" {
  instance           = lxd_instance.instance1.name
  source_path        = "../acctest/fixtures/test-file.txt"
  target_path        = "/foo/bar.txt"
  create_directories = true
}
	`, name, acctest.TestImage)
}

func testAccInstanceFile_project(project, instance string) string {
	return fmt.Sprintf(`
resource "lxd_project" "project1" {
  name = "%s"
  config = {
    "features.images"   = false
    "features.profiles" = false
  }
}

resource "lxd_instance" "instance1" {
  name    = "%s"
  image   = "%s"
  project = lxd_project.project1.name
}

resource "lxd_instance_file" "file1" {
  instance           = lxd_instance.instance1.name
  project   	     = lxd_project.project1.name
  source_path        = "../acctest/fixtures/test-file.txt"
  target_path        = "/foo/bar.txt"
  create_directories = true
}
	`, project, instance, acctest.TestImage)
}
