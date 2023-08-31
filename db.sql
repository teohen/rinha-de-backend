CREATE TEXT SEARCH CONFIGURATION BUSCA (COPY = portuguese);
ALTER TEXT SEARCH CONFIGURATION BUSCA ALTER MAPPING FOR hword, hword_part, word WITH portuguese_stem;

ALTER DATABASE rinhadebackend SET synchronous_commit=OFF;

ALTER SYSTEM SET max_connections = 1000;

ALTER SYSTEM SET shared_buffers TO "425MB";


CREATE TABLE IF NOT EXISTS pessoas (
  id UUID PRIMARY KEY NOT NULL,
  apelido VARCHAR(32) UNIQUE NOT NULL,
  nome VARCHAR(100),
  nascimento CHAR(10) NOT NULL,
  stack TEXT NULL
  BUSCA_TRGM TEXT GENERATED ALWAYS AS (
  LOWER(nome) || LOWER(apelido) || LOWER(stack)
  ) STORED
  );

CREATE EXTENSION PG_TRGM;

CREATE INDEX CONCURRENTLY IF NOT EXISTS IDX_PESSOAS_BUSCA_TGRM ON pessoas USING GIST (BUSCA_TRGM GIST_TRGM_OPS);
CREATE INDEX CONCURRENTLY IF NOT EXISTS IDXPESSOAS_GIN ON pessoas USING GIN (BUSCA_TRGM gin_trgm_ops);
