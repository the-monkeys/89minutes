# Run the following command to run opensearch container
`sudo docker run -d -p 9200:9200 -p 9600:9600 -e "discovery.type=single-node" opensearchproject/opensearch:latest`

The above command starts a container in the background and we can have a local environment to test data against the database.