package main_test

import (
	"github.com/cloud66-oss/starter/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
)

func runStarter() (string, error) {
	command := exec.Command(binPath)
	command_out, err := command.Output()
	output := string(command_out)
	return output, err
}

func runStarterWithProjectNoTemplates(projectFixture string) (string, error) {
	command := exec.Command(binPath, "-y", "-p", projectFixture+"/src", "-g", "dockerfile, c66-service")
	command_out, err := command.Output()
	output := string(command_out)
	return output, err
}

func runStarterWithProject(projectFixture string) (string, error) {
	command := exec.Command(binPath, "-y", "-p", projectFixture+"/src", "-templates", "templates/", "-g", "dockerfile, c66-service")
	command_out, err := command.Output()
	output := string(command_out)
	return output, err
}

func runStarterWithProjectGeneratingOnlyDockerfile(projectFixture string) (string, error) {
	command := exec.Command(binPath, "-y", "-p", projectFixture+"/src", "-templates", "templates/")
	command_out, err := command.Output()
	output := string(command_out)
	return output, err
}

func runStarterWithProjectGeneratingKubernetesYmlFromServiceYml(projectFixturePath string) (string, error) {
	command := exec.Command(binPath, "-y", "-p", projectFixturePath+"/src", "-g", "kubernetes")
	command_out, err := command.Output()
	output := string(command_out)
	return output, err
}

func cleanupGeneratedFiles(projectFixture string) {
	os.Remove(projectFixture + "/src/Dockerfile")
	os.Remove(projectFixture + "/src/service.yml")
}

// NOTE: starter will be detected as the test projects git repo, so in order
// for tests to always work we replace the current starter branch (which may
// change) to 'master' in the generated file.
func convertServiceYaml(generated []byte) []byte {
	generated = regexp.MustCompile(`git_branch: .*`).ReplaceAll(generated, []byte("git_branch: master"))
	generated = regexp.MustCompile(`git_url: .*`).ReplaceAll(generated, []byte("git_url: git@github.com:cloud66/starter.git"))
	return generated
}

/*
var _ = Describe("Start Starter without flags", func() {
	It("should run starter but fail because no project", func() {
		output, err := runStarter()
		outputLines := strings.Split(output, "\n")
		message := outputLines[len(outputLines) - 2]
		Expect(strings.Contains(message, "Failed to detect framework due to: Could not detect any of the supported frameworks")).To(BeTrue())
		Expect(err).NotTo(HaveOccurred())
	})
})
*/

/*
var _ = Describe("Generating all files with Starter", func() {
	Context("using a NodeJS Express project with no database but not in GIT", func() {
		var projectFixturePath string = "/usr/local/go/src/github.com/cloud66"
		BeforeEach(func() {
			_, err := runStarterWithProjectNoTemplates(projectFixturePath)
			Expect(err).NotTo(HaveOccurred())
		})

    	AfterEach(func() {
			cleanupGeneratedFiles(projectFixturePath)
		})

		It("should generate a Dockerfile", func() {
			dockerfile_expected,_ := ioutil.ReadFile(projectFixturePath + "/expected/Dockerfile")
			dockerfile_generated,_ := ioutil.ReadFile(projectFixturePath + "/src/Dockerfile")
			Expect(dockerfile_generated).To(Equal(dockerfile_expected))
		})

		It("should generate a service.yml", func() {
			service_yaml_expected,_ := ioutil.ReadFile(projectFixturePath + "/expected/service.yml")
			service_yaml_generated,_ := ioutil.ReadFile(projectFixturePath + "/src/service.yml")
			service_yaml_generated = convertServiceYaml(service_yaml_generated)
			Expect(service_yaml_generated).To(Equal(service_yaml_expected))
		})
	})
})*/

var _ = Describe("Generating all files with Starter", func() {
	Context("using a PHP Laravel project with a mysql database", func() {
		var projectFixturePath string = "test/php/laravel_bare_mysql"
		BeforeEach(func() {
			_, err := runStarterWithProject(projectFixturePath)
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			cleanupGeneratedFiles(projectFixturePath)
		})

		It("should generate a Dockerfile", func() {
			dockerfile_expected, _ := ioutil.ReadFile(projectFixturePath + "/expected/Dockerfile")
			dockerfile_generated, _ := ioutil.ReadFile(projectFixturePath + "/src/Dockerfile")
			Expect(dockerfile_generated).To(Equal(dockerfile_expected))
		})

		It("should generate a service.yml", func() {
			service_yaml_expected, _ := ioutil.ReadFile(projectFixturePath + "/expected/service.yml")
			service_yaml_generated, _ := ioutil.ReadFile(projectFixturePath + "/src/service.yml")
			service_yaml_generated = convertServiceYaml(service_yaml_generated)
			Expect(service_yaml_generated).To(Equal(service_yaml_expected))
		})
	})

	/*	Context("using a NodeJS Express project with rethinkdb", func() {
				var projectFixturePath string = "test/node/express_rethinkdb"
				BeforeEach(func() {
					_, err := runStarterWithProject(projectFixturePath)
					Expect(err).NotTo(HaveOccurred())
				})

		    	AfterEach(func() {
					cleanupGeneratedFiles(projectFixturePath)
				})

				It("should generate a Dockerfile", func() {
					dockerfile_expected,_ := ioutil.ReadFile(projectFixturePath + "/expected/Dockerfile")
					dockerfile_generated,_ := ioutil.ReadFile(projectFixturePath + "/src/Dockerfile")
					Expect(dockerfile_generated).To(Equal(dockerfile_expected))
				})

				It("should generate a service.yml", func() {
					service_yaml_expected,_ := ioutil.ReadFile(projectFixturePath + "/expected/service.yml")
					service_yaml_generated,_ := ioutil.ReadFile(projectFixturePath + "/src/service.yml")
					service_yaml_generated = convertServiceYaml(service_yaml_generated)
					Expect(service_yaml_generated).To(Equal(service_yaml_expected))
				})

			})
	*/

	Context("using a NodeJS Express project with no database", func() {
		var projectFixturePath string = "test/node/express_bare"
		BeforeEach(func() {
			_, err := runStarterWithProject(projectFixturePath)
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			cleanupGeneratedFiles(projectFixturePath)
		})

		It("should generate a Dockerfile", func() {
			dockerfile_expected, _ := ioutil.ReadFile(projectFixturePath + "/expected/Dockerfile")
			dockerfile_generated, _ := ioutil.ReadFile(projectFixturePath + "/src/Dockerfile")
			Expect(dockerfile_generated).To(Equal(dockerfile_expected))
		})

		It("should generate a service.yml", func() {
			service_yaml_expected, _ := ioutil.ReadFile(projectFixturePath + "/expected/service.yml")
			service_yaml_generated, _ := ioutil.ReadFile(projectFixturePath + "/src/service.yml")
			service_yaml_generated = convertServiceYaml(service_yaml_generated)
			Expect(service_yaml_generated).To(Equal(service_yaml_expected))
		})
	})
	Context("using a NodeJS Express project with a mysql database", func() {
		var projectFixturePath string = "test/node/express_bare_mysql"
		BeforeEach(func() {
			_, err := runStarterWithProject(projectFixturePath)
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			cleanupGeneratedFiles(projectFixturePath)
		})

		It("should generate a Dockerfile", func() {
			dockerfile_expected, _ := ioutil.ReadFile(projectFixturePath + "/expected/Dockerfile")
			dockerfile_generated, _ := ioutil.ReadFile(projectFixturePath + "/src/Dockerfile")
			Expect(dockerfile_generated).To(Equal(dockerfile_expected))
		})

		It("should generate a service.yml", func() {
			service_yaml_expected, _ := ioutil.ReadFile(projectFixturePath + "/expected/service.yml")
			service_yaml_generated, _ := ioutil.ReadFile(projectFixturePath + "/src/service.yml")
			service_yaml_generated = convertServiceYaml(service_yaml_generated)
			Expect(service_yaml_generated).To(Equal(service_yaml_expected))
		})
	})
	Context("using a NodeJS Express project with a mongo database", func() {
		var projectFixturePath string = "test/node/express_bare_mongo"
		BeforeEach(func() {
			_, err := runStarterWithProject(projectFixturePath)
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			cleanupGeneratedFiles(projectFixturePath)
		})

		It("should generate a Dockerfile", func() {
			dockerfile_expected, _ := ioutil.ReadFile(projectFixturePath + "/expected/Dockerfile")
			dockerfile_generated, _ := ioutil.ReadFile(projectFixturePath + "/src/Dockerfile")
			Expect(dockerfile_generated).To(Equal(dockerfile_expected))
		})

		It("should generate a service.yml", func() {
			service_yaml_expected, _ := ioutil.ReadFile(projectFixturePath + "/expected/service.yml")
			service_yaml_generated, _ := ioutil.ReadFile(projectFixturePath + "/src/service.yml")
			service_yaml_generated = convertServiceYaml(service_yaml_generated)
			Expect(service_yaml_generated).To(Equal(service_yaml_expected))
		})
	})
	Context("using a NodeJS Express project with Procfile and no database", func() {
		var projectFixturePath string = "test/node/express_procfile"
		BeforeEach(func() {
			_, err := runStarterWithProject(projectFixturePath)
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			cleanupGeneratedFiles(projectFixturePath)
		})

		It("should generate a Dockerfile", func() {
			dockerfile_expected, _ := ioutil.ReadFile(projectFixturePath + "/expected/Dockerfile")
			dockerfile_generated, _ := ioutil.ReadFile(projectFixturePath + "/src/Dockerfile")
			Expect(dockerfile_generated).To(Equal(dockerfile_expected))
		})

		It("should generate a service.yml", func() {
			service_yaml_expected, _ := ioutil.ReadFile(projectFixturePath + "/expected/service.yml")
			service_yaml_generated, _ := ioutil.ReadFile(projectFixturePath + "/src/service.yml")
			service_yaml_generated = convertServiceYaml(service_yaml_generated)
			Expect(service_yaml_generated).To(Equal(service_yaml_expected))
		})
	})
	Context("using a Rails project with a Mysql database", func() {
		var projectFixturePath string = "test/ruby/rails_mysql"

		BeforeEach(func() {
			_, err := runStarterWithProject(projectFixturePath)
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			cleanupGeneratedFiles(projectFixturePath)
		})

		It("should generate a Dockerfile", func() {
			dockerfile_expected, _ := ioutil.ReadFile(projectFixturePath + "/expected/Dockerfile")
			dockerfile_generated, _ := ioutil.ReadFile(projectFixturePath + "/src/Dockerfile")
			Expect(dockerfile_generated).To(Equal(dockerfile_expected))
		})

		It("should generate a service.yml", func() {
			service_yaml_expected, _ := ioutil.ReadFile(projectFixturePath + "/expected/service.yml")
			service_yaml_generated, _ := ioutil.ReadFile(projectFixturePath + "/src/service.yml")
			service_yaml_generated = convertServiceYaml(service_yaml_generated)
			Expect(service_yaml_generated).To(Equal(service_yaml_expected))
		})
	})

	Context("using a Rails project running Unicorn and using a Mysql database", func() {
		var projectFixturePath string = "test/ruby/rails_unicorn_mysql"

		BeforeEach(func() {
			_, err := runStarterWithProject(projectFixturePath)
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			cleanupGeneratedFiles(projectFixturePath)
		})

		It("should generate a Dockerfile", func() {
			dockerfile_expected, _ := ioutil.ReadFile(projectFixturePath + "/expected/Dockerfile")
			dockerfile_generated, _ := ioutil.ReadFile(projectFixturePath + "/src/Dockerfile")
			Expect(dockerfile_generated).To(Equal(dockerfile_expected))
		})

		It("should generate a service.yml", func() {
			service_yaml_expected, _ := ioutil.ReadFile(projectFixturePath + "/expected/service.yml")
			service_yaml_generated, _ := ioutil.ReadFile(projectFixturePath + "/src/service.yml")
			service_yaml_generated = convertServiceYaml(service_yaml_generated)
			Expect(service_yaml_generated).To(Equal(service_yaml_expected))
		})
	})

	Context("using a Rails project running Unicorn, some workers and using a Redis and Postgresql database", func() {
		var projectFixturePath string = "test/ruby/rails_jobs_unicorn_redis_postgresql"

		BeforeEach(func() {
			_, err := runStarterWithProject(projectFixturePath)
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			cleanupGeneratedFiles(projectFixturePath)
		})

		It("should generate a Dockerfile", func() {
			dockerfile_expected, _ := ioutil.ReadFile(projectFixturePath + "/expected/Dockerfile")
			dockerfile_generated, _ := ioutil.ReadFile(projectFixturePath + "/src/Dockerfile")
			Expect(dockerfile_generated).To(Equal(dockerfile_expected))
		})

		It("should generate a service.yml", func() {
			service_yaml_expected, _ := ioutil.ReadFile(projectFixturePath + "/expected/service.yml")
			service_yaml_generated, _ := ioutil.ReadFile(projectFixturePath + "/src/service.yml")
			service_yaml_generated = convertServiceYaml(service_yaml_generated)
			Expect(service_yaml_generated).To(Equal(service_yaml_expected))
		})
	})
	Context("using a NodeJS Express project with a pg database", func() {
		var projectFixturePath string = "test/node/express_pg_procfile"
		BeforeEach(func() {
			_, err := runStarterWithProject(projectFixturePath)
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			cleanupGeneratedFiles(projectFixturePath)
		})

		It("should generate a Dockerfile", func() {
			dockerfile_expected, _ := ioutil.ReadFile(projectFixturePath + "/expected/Dockerfile")
			dockerfile_generated, _ := ioutil.ReadFile(projectFixturePath + "/src/Dockerfile")
			Expect(dockerfile_generated).To(Equal(dockerfile_expected))
		})

		It("should generate a service.yml", func() {
			service_yaml_expected, _ := ioutil.ReadFile(projectFixturePath + "/expected/service.yml")
			service_yaml_generated, _ := ioutil.ReadFile(projectFixturePath + "/src/service.yml")
			service_yaml_generated = convertServiceYaml(service_yaml_generated)
			Expect(service_yaml_generated).To(Equal(service_yaml_expected))
		})
	})
	Context("using a NodeJS Express project with a mongodb database", func() {
		var projectFixturePath string = "test/node/express_mongodb"
		BeforeEach(func() {
			_, err := runStarterWithProject(projectFixturePath)
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			cleanupGeneratedFiles(projectFixturePath)
		})

		It("should generate a Dockerfile", func() {
			dockerfile_expected, _ := ioutil.ReadFile(projectFixturePath + "/expected/Dockerfile")
			dockerfile_generated, _ := ioutil.ReadFile(projectFixturePath + "/src/Dockerfile")
			Expect(dockerfile_generated).To(Equal(dockerfile_expected))
		})

		It("should generate a service.yml", func() {
			service_yaml_expected, _ := ioutil.ReadFile(projectFixturePath + "/expected/service.yml")
			service_yaml_generated, _ := ioutil.ReadFile(projectFixturePath + "/src/service.yml")
			service_yaml_generated = convertServiceYaml(service_yaml_generated)
			Expect(service_yaml_generated).To(Equal(service_yaml_expected))
		})
	})
	Context("using a NodeJS Express project with a mongodb database", func() {
		var projectFixturePath string = "test/node/express_mongodb_procfile"
		BeforeEach(func() {
			_, err := runStarterWithProject(projectFixturePath)
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			cleanupGeneratedFiles(projectFixturePath)
		})

		It("should generate a Dockerfile", func() {
			dockerfile_expected, _ := ioutil.ReadFile(projectFixturePath + "/expected/Dockerfile")
			dockerfile_generated, _ := ioutil.ReadFile(projectFixturePath + "/src/Dockerfile")
			Expect(dockerfile_generated).To(Equal(dockerfile_expected))
		})

		It("should generate a service.yml", func() {
			service_yaml_expected, _ := ioutil.ReadFile(projectFixturePath + "/expected/service.yml")
			service_yaml_generated, _ := ioutil.ReadFile(projectFixturePath + "/src/service.yml")
			service_yaml_generated = convertServiceYaml(service_yaml_generated)
			Expect(service_yaml_generated).To(Equal(service_yaml_expected))
		})
	})
	Context("using a NodeJS Express project with a mongodb and redis database", func() {
		var projectFixturePath string = "test/node/express_mongodb_redis"
		BeforeEach(func() {
			_, err := runStarterWithProject(projectFixturePath)
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			cleanupGeneratedFiles(projectFixturePath)
		})

		It("should generate a Dockerfile", func() {
			dockerfile_expected, _ := ioutil.ReadFile(projectFixturePath + "/expected/Dockerfile")
			dockerfile_generated, _ := ioutil.ReadFile(projectFixturePath + "/src/Dockerfile")
			Expect(dockerfile_generated).To(Equal(dockerfile_expected))
		})

		It("should generate a service.yml", func() {
			service_yaml_expected, _ := ioutil.ReadFile(projectFixturePath + "/expected/service.yml")
			service_yaml_generated, _ := ioutil.ReadFile(projectFixturePath + "/src/service.yml")
			service_yaml_generated = convertServiceYaml(service_yaml_generated)
			Expect(service_yaml_generated).To(Equal(service_yaml_expected))
		})
	})
	Context("using a NodeJS Express project with a mysql and postgres database", func() {
		var projectFixturePath string = "test/node/express_mysql_pg"
		BeforeEach(func() {
			_, err := runStarterWithProject(projectFixturePath)
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			cleanupGeneratedFiles(projectFixturePath)
		})

		It("should generate a Dockerfile", func() {
			dockerfile_expected, _ := ioutil.ReadFile(projectFixturePath + "/expected/Dockerfile")
			dockerfile_generated, _ := ioutil.ReadFile(projectFixturePath + "/src/Dockerfile")
			Expect(dockerfile_generated).To(Equal(dockerfile_expected))
		})

		It("should generate a service.yml", func() {
			service_yaml_expected, _ := ioutil.ReadFile(projectFixturePath + "/expected/service.yml")
			service_yaml_generated, _ := ioutil.ReadFile(projectFixturePath + "/src/service.yml")
			service_yaml_generated = convertServiceYaml(service_yaml_generated)
			Expect(service_yaml_generated).To(Equal(service_yaml_expected))
		})
	})
})
var _ = Describe("Generating only Dockerfile with Starter", func() {
	Context("using a Rails project with a Mysql database", func() {
		var projectFixturePath string = "test/ruby/rails_mysql"

		BeforeEach(func() {
			_, err := runStarterWithProjectGeneratingOnlyDockerfile(projectFixturePath)
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			cleanupGeneratedFiles(projectFixturePath)
		})

		It("should generate a Dockerfile", func() {
			dockerfile_expected, _ := ioutil.ReadFile(projectFixturePath + "/expected/Dockerfile")
			dockerfile_generated, _ := ioutil.ReadFile(projectFixturePath + "/src/Dockerfile")
			Expect(dockerfile_generated).To(Equal(dockerfile_expected))
		})

		It("should not generate a service.yml", func() {
			Expect(common.FileExists(projectFixturePath + "/src/service.yml")).To(BeFalse())
		})
	})

})

var _ = Describe("Generating kubernetes.yml from service.yml", func() {
	Context("using an empty project - service.yml populated with all possible keys", func() {
		var projectFixturePath string = "test/service_yml/service_to_kubes_1"

		BeforeEach(func() {
			_, err := runStarterWithProjectGeneratingKubernetesYmlFromServiceYml(projectFixturePath)
			Expect(err).NotTo(HaveOccurred())
		})
		AfterEach(func() {
			os.Remove(projectFixturePath + "/src/kubernetes.yml")
		})

		It("should generate kubernetes.yml", func() {
			kubernetes_expected, _ := ioutil.ReadFile(projectFixturePath + "/expected/kubernetes.yml")
			kubernetes_generated, _ := ioutil.ReadFile(projectFixturePath + "/src/kubernetes.yml")
			Expect(kubernetes_generated).To(Equal(kubernetes_expected))
		})
	})
	Context("using an empty project - functional service.yml", func() {
		var projectFixturePath string = "test/service_yml/service_to_kubes_2"

		BeforeEach(func() {
			_, err := runStarterWithProjectGeneratingKubernetesYmlFromServiceYml(projectFixturePath)
			Expect(err).NotTo(HaveOccurred())
		})
		AfterEach(func() {
			os.Remove(projectFixturePath + "/src/kubernetes.yml")
		})

		It("should generate kubernetes.yml", func() {
			kubernetes_expected, _ := ioutil.ReadFile(projectFixturePath + "/expected/kubernetes.yml")
			kubernetes_generated, _ := ioutil.ReadFile(projectFixturePath + "/src/kubernetes.yml")
			Expect(kubernetes_generated).To(Equal(kubernetes_expected))
		})
	})
})
