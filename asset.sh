echo "drop table asset" | sudo -u postgres psql -U postgres postgres
#id | symbol | price_usd | balance | paid_at | created_at | updated_at 
# TODO: need to infer balance, created_at, updated_at from other tables
echo "select asset_id as id, symbol, price_usd into table asset from assets;" | sudo -u postgres psql -U postgres postgres
sudo -u postgres pg_dump -t asset > asset.sql
sed -i -e 's/postgres/setflags/g' asset.sql
source secrets/db.sh
echo "drop table asset;" | psql
psql < asset.sql
