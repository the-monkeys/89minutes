FROM opensearchproject/opensearch:2.3.0
RUN /usr/share/opensearch/bin/opensearch-plugin install --batch opensearch-security
