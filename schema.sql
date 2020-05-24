CREATE TABLE IF NOT EXISTS flags (
  flag_id           VARCHAR(36) PRIMARY KEY CHECK (flag_id ~* '^[0-9a-f-]{36,36}$'),
  payer_id          VARCHAR(36) NOT NULL CHECK (payer_id ~* '^[0-9a-f-]{36,36}$'),
  task              VARCHAR(256) NOT NULL,
  days	    	    BIGINT NOT NULL DEFAULT 1,
  asset_id          VARCHAR(36) NOT NULL CHECK (asset_id ~* '^[0-9a-f-]{36,36}$'),
  amount            VARCHAR(128) NOT NULL,
  remaining_amount  VARCHAR(128) NOT NULL,
  remaining_days    BIGINT NOT NULL,
  times_achieved    BIGINT NOT NULL,
  max_witness       BIGINT NOT NULL,
  state             VARCHAR(36) NOT NULL,
  created_at        TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS flags_state_createdx ON flags(state, created_at);

CREATE TABLE IF NOT EXISTS evidences (
  attachement_id        VARCHAR(36) PRIMARY KEY CHECK (attachement_id ~* '^[0-9a-f-]{36,36}$'),
  flag_id               VARCHAR(36) NOT NULL CHECK (flag_id ~* '^[0-9a-f-]{36,36}$'),
  file                  TEXT NOT NULL,
  type              	VARCHAR(512) NOT NULL,
  created_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS evidence_datesx ON evidences(created_at);

CREATE TABLE IF NOT EXISTS witnesses (
  flag_id           VARCHAR(36) NOT NULL REFERENCES flags(flag_id) ON DELETE CASCADE,
  payee_id           VARCHAR(36) NOT NULL CHECK (payee_id ~* '^[0-9a-f-]{36,36}$'),
  created_at        TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  amount            VARCHAR(128) NOT NULL,
  paid_at           TIMESTAMP WITH TIME ZONE,
  verified	    SMALLINT NOT NULL DEFAULT 0,
  PRIMARY KEY(flag_id, payee_id)
);

CREATE INDEX IF NOT EXISTS witnesses_created_paidx ON witnesses(created_at);

CREATE TABLE IF NOT EXISTS users (
  user_id           VARCHAR(36) PRIMARY KEY CHECK (user_id ~* '^[0-9a-f-]{36,36}$'),
  identity_number   BIGINT NOT NULL,
  full_name         VARCHAR(512) NOT NULL DEFAULT '',
  access_token      VARCHAR(512) NOT NULL DEFAULT '',
  avatar_url        VARCHAR(1024) NOT NULL DEFAULT '',
  trace_id          VARCHAR(36) NOT NULL CHECK (trace_id ~* '^[0-9a-f-]{36,36}$'),
  state             VARCHAR(128) NOT NULL,
  locale 	    VARCHAR(2) NOT NULL DEFAULT 'zh'
);

CREATE UNIQUE INDEX IF NOT EXISTS users_identityx ON users(identity_number);

CREATE TABLE IF NOT EXISTS assets (
  asset_id         VARCHAR(36) PRIMARY KEY CHECK (asset_id ~* '^[0-9a-f-]{36,36}$'),
  symbol           VARCHAR(512) NOT NULL,
  price_usd        VARCHAR(128) NOT NULL,
  balance	   VARCHAR(128) NOT NULL
);

CREATE INDEX IF NOT EXISTS assets_symbol ON assets(symbol);
