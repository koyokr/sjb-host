# sjb-host

[![travis-ci](https://api.travis-ci.org/koyokr/sjb-host.svg?branch=master)](https://travis-ci.org/koyokr/sjb-host)

```sql
CREATE TABLE domains (
id serial PRIMARY KEY,
host varchar(255) NOT NULL,
round_robin bool NOT NULL,
has_blocked bool NOT NULL,
UNIQUE(host)
);

CREATE TABLE ipss (
id serial PRIMARY KEY,
value text NOT NULL,
unique(value)
);

CREATE TABLE domain_to_ipss (
domain_id int not null REFERENCES domains,
ips_id int not null REFERENCES ipss
);

CREATE INDEX domain_to_ipss_domain_id_key ON domain_to_ipss (domain_id);
```
