precommit/init:
	brew install pre-commit
	pre-commit install

ci_test:
	go test ./... -v -coverprofile .coverage.txt

ci_sonarqube:
	sonar-scanner \
        -Dsonar.host.url=#SONARQUBE_URL# \
        -Dsonar.projectVersion=#APP_VERSION# \
        -Dsonar.go.coverage.reportPaths=reports/.coverage.txt \

ci_sonarqube_local:
	sonar-scanner
	# See configurations in sonar-project.properties

coverage_scanner:
	$(MAKE) ci_test
	$(MAKE) ci_sonarqube_local
