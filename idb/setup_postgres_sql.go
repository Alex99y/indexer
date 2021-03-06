// GENERATED CODE from source setup_postgres.sql via go generate

package idb

const setup_postgres_sql = `-- This file is setup_postgres.sql which gets compiled into go source using a go:generate statement in postgres.go
--
-- TODO? replace all 'addr bytea' with 'addr_id bigint' and a mapping table? makes addrs an 8 byte int that fits in a register instead of a 32 byte string

CREATE TABLE IF NOT EXISTS protocol (
  version text PRIMARY KEY,
  proto jsonb
);

CREATE TABLE IF NOT EXISTS block_header (
round bigint PRIMARY KEY,
realtime timestamp without time zone NOT NULL,
rewardslevel bigint NOT NULL,
header jsonb NOT NULL
);
CREATE INDEX IF NOT EXISTS block_header_time ON block_header (realtime);

CREATE TABLE IF NOT EXISTS txn (
round bigint NOT NULL,
intra smallint NOT NULL,
typeenum smallint NOT NULL,
asset bigint NOT NULL, -- 0=Algos, otherwise AssetIndex
txid bytea NOT NULL, -- [32]byte
txnbytes bytea NOT NULL,
txn jsonb NOT NULL,
extra jsonb,
PRIMARY KEY ( round, intra )
);

-- NOT a unique index because we don't guarantee txid is unique outside of its 1000 rounds.
CREATE INDEX IF NOT EXISTS txn_by_tixid ON txn ( txid );

CREATE TABLE IF NOT EXISTS txn_participation (
addr bytea NOT NULL,
round bigint NOT NULL,
intra smallint NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS txn_participation_i ON txn_participation ( addr, round DESC, intra DESC );

-- bookeeping for local file import
CREATE TABLE IF NOT EXISTS imported (path text);

-- expand data.basics.AccountData
CREATE TABLE IF NOT EXISTS account (
  addr bytea primary key,
  microalgos bigint NOT NULL, -- okay because less than 2^54 Algos
  rewardsbase bigint NOT NULL,
  keytype varchar(8), -- sig,msig,lsig
  account_data jsonb -- data.basics.AccountData except AssetParams and Assets and MicroAlgos and RewardsBase
);

-- data.basics.AccountData Assets[asset id] AssetHolding{}
CREATE TABLE IF NOT EXISTS account_asset (
  addr bytea NOT NULL, -- [32]byte
  assetid bigint NOT NULL,
  amount numeric(20) NOT NULL, -- need the full 18446744073709551615
  frozen boolean NOT NULL,
  PRIMARY KEY (addr, assetid)
);

-- data.basics.AccountData AssetParams[index] AssetParams{}
CREATE TABLE IF NOT EXISTS asset (
  index bigint PRIMARY KEY,
  creator_addr bytea NOT NULL,
  params jsonb NOT NULL -- data.basics.AssetParams -- TODO index some fields?
);
-- TODO: index on creator_addr?

-- subsumes ledger/accountdb.go accounttotals and acctrounds
-- "state":{online, onlinerewardunits, offline, offlinerewardunits, notparticipating, notparticipatingrewardunits, rewardslevel, round bigint}
CREATE TABLE IF NOT EXISTS metastate (
  k text primary key,
  v jsonb
);
`
