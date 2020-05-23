all: docs/models.png README.md

docs/models.png: docs/models.puml
	cp out/docs/models/models.png docs/

README.md: docs/models.yml
	markdown-swagger ./docs/models.yml ./README.md
