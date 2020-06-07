echo "drop table flag" | sudo -u postgres psql -U postgres postgres
echo "select packet_id as id, packets.user_id as payer_id, greeting as task, period as days, total_count as max_witness, asset_id, amount, period - remaining_count as times_achieved, packets.state as status, remaining_count as remaining_days, remaining_amount, created_at, created_at as updated_at, u.full_name as payer_name into table flag from packets, users u where packets.user_id = u.user_id;" | sudo -u postgres psql -U postgres postgres
sudo -u postgres pg_dump -t flag > flag.sql
sed -i -e 's/postgres/setflags/g' flag.sql
source secrets/db.sh
echo "drop table flag;" | psql
psql < flag.sql
