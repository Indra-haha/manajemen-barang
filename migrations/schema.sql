CREATE TABLE IF NOT EXISTS barang (
  id         BIGSERIAL    PRIMARY KEY,
  kode       VARCHAR(20)  NOT NULL UNIQUE,
  nama       VARCHAR(100) NOT NULL,
  kategori   VARCHAR(50)  NOT NULL DEFAULT 'umum',
  stok       INTEGER      NOT NULL DEFAULT 0 CHECK (stok >= 0),
  harga      BIGINT       NOT NULL DEFAULT 0 CHECK (harga >= 0),
  created_at TIMESTAMPTZ  NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_barang_kategori ON barang (kategori);

CREATE TABLE IF NOT EXISTS users (
  id            BIGSERIAL    PRIMARY KEY,
  username      VARCHAR(50)  NOT NULL UNIQUE,
  password_hash VARCHAR(100) NOT NULL,
  created_at    TIMESTAMPTZ  NOT NULL DEFAULT now()
);
