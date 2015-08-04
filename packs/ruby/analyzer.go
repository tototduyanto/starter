package ruby

import (
	"fmt"
	"path/filepath"

	"github.com/cloud66/starter/common"
	"github.com/cloud66/starter/packs"
	"github.com/cloud66/starter/packs/ruby/webservers"
)

type Analyzer struct {
	packs.AnalyzerBase
	Gemfile string
}

func (a *Analyzer) Analyze() (*Analysis, error) {
	a.Gemfile = filepath.Join(a.RootDir, "Gemfile")
	gitURL, gitBranch, buildRoot, err := a.ProjectMetadata()
	if err != nil {
		return nil, err
	}

	packages := a.GuessPackages()
	version := a.FindVersion()
	dbs := a.ConfirmDatabases(a.FindDatabases())
	envVars := a.EnvVars()

	services, err := a.AnalyzeServices(a, envVars, gitBranch, gitURL, buildRoot)
	if err != nil {
		return nil, err
	}

	analysis := &Analysis{
		AnalysisBase: packs.AnalysisBase{
			PackName:  a.GetPack().Name(),
			GitBranch: gitBranch,
			GitURL:    gitURL,
			Messages:  a.Messages},
		ServiceYAMLContext: &ServiceYAMLContext{packs.ServiceYAMLContextBase{Services: services, Dbs: dbs.Items}},
		DockerfileContext:  &DockerfileContext{packs.DockerfileContextBase{Version: version, Packages: packages}}}
	return analysis, nil
}

func (a *Analyzer) FillServices(services *[]*common.Service) error {
	service := a.GetOrCreateWebService(services)
	service.Ports = []*common.PortMapping{common.NewPortMapping()}
	hasFoundServer, server := a.detectWebServer(service.Command)

	if service.Command == "" {
		isRails, _ := common.GetGemVersion(a.Gemfile, "rails")
		if isRails {
			service.Command = "bundle exec rails s _env:RAILS_ENV"
			service.Ports[0].Container = "3000"
		} else {
			service.Command = "bundle exec rackup s _env:RACK_ENV"
			service.Ports[0].Container = "9292"
		}
	} else {
		if hasFoundServer {
			service.Ports[0].Container = server.Port(service.Command)
		} else {
			hasFound, port := common.ParsePort(service.Command)
			if hasFound {
				service.Ports[0].Container = port
			} else {
				if !a.ShouldPrompt {
					return fmt.Errorf("Could not find port to open corresponding to command '%s'", service.Command)
				}
				service.Ports[0].Container = common.AskUser(fmt.Sprintf("Which port to open to run web service with command '%s'?", service.Command))
			}
		}
	}

	service.BuildCommand = a.AskForCommand("bundle exec rake db:schema:load", "build")
	service.DeployCommand = a.AskForCommand("bundle exec rake db:migrate", "deployment")

	return nil
}

func (a *Analyzer) HasPackage(pack string) bool {
	hasFound, _ := common.GetGemVersion(a.Gemfile, pack)
	return hasFound
}

func (a *Analyzer) detectWebServer(command string) (hasFound bool, server packs.WebServer) {
	unicorn := &webservers.Unicorn{}
	thin := &webservers.Thin{}
	servers := []packs.WebServer{unicorn, thin}
	return a.AnalyzerBase.DetectWebServer(a, command, servers)
}

func (a *Analyzer) GuessPackages() *common.Lister {
	packages := common.NewLister()
	if hasRmagick, _ := common.GetGemVersion(a.Gemfile, "rmagick"); hasRmagick {
		fmt.Println(common.MsgL2, "----> Found Image Magick", common.MsgReset)
		packages.Add("imagemagick", "libmagickwand-dev")
	}

	if hasSqlite, _ := common.GetGemVersion(a.Gemfile, "sqlite"); hasSqlite {
		packages.Add("libsqlite3-dev")
		fmt.Println(common.MsgL2, "----> Found sqlite", common.MsgReset)
	}
	return packages
}

func (a *Analyzer) FindVersion() string {
	foundRuby, rubyVersion := common.GetRubyVersion(a.Gemfile)
	return a.ConfirmVersion(foundRuby, rubyVersion, "latest")
}

func (a *Analyzer) FindDatabases() *common.Lister {
	dbs := common.NewLister()
	if hasMysql, _ := common.GetGemVersion(a.Gemfile, "mysql2"); hasMysql {
		dbs.Add("mysql")
	}

	if hasPg, _ := common.GetGemVersion(a.Gemfile, "pg"); hasPg {
		dbs.Add("postgresql")
	}

	if hasRedis, _ := common.GetGemVersion(a.Gemfile, "redis"); hasRedis {
		dbs.Add("redis")
	}

	if hasMongoDB, _ := common.GetGemVersion(a.Gemfile, "mongo", "mongo_mapper", "dm-mongo-adapter", "mongoid"); hasMongoDB {
		dbs.Add("mongodb")
	}

	if hasElasticsearch, _ := common.GetGemVersion(a.Gemfile, "elasticsearch", "tire", "flex", "chewy"); hasElasticsearch {
		dbs.Add("elasticsearch")
	}

	if hasDatabaseYaml := common.FileExists("config/database.yml"); hasDatabaseYaml {
		fmt.Println(common.MsgL2, "----> Found config/database.yml", common.MsgReset)
		a.Messages.Add(
			fmt.Sprintf("%s %s-> %s",
				"database.yml: Make sure you are using environment variables.",
				common.MsgReset, "http://help.cloud66.com/deployment/environment-variables"))
	}

	if hasMongoIdYaml := common.FileExists("config/mongoid.yml"); hasMongoIdYaml {
		fmt.Println(common.MsgL2, "----> Found config/mongoid.yml", common.MsgReset)
		a.Messages.Add(
			fmt.Sprintf("%s %s-> %s",
				"mongoid.yml: Make sure you are using environment variables.",
				common.MsgReset, "http://help.cloud66.com/deployment/environment-variables"))
	}
	return dbs
}

func (a *Analyzer) EnvVars() []*common.EnvVar {
	return []*common.EnvVar{
		&common.EnvVar{Key: "RAILS_ENV", Value: a.Environment},
		&common.EnvVar{Key: "RACK_ENV", Value: a.Environment}}
}
