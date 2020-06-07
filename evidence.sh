echo "drop table evidence" | sudo -u postgres psql -U postgres postgres
# id | attachment_id | flag_id | url | type | created_at | updated_at 
# TODO: need to infer attachment_id, url from Mixin Messenger API
echo "select m.message_id as id, p.packet_id as flag_id, p.user_id as payee_id, p.done as state, m.category, m.created_at, m.updated_at as type into table evidence from participants p, messages m where p.user_id = m.user_id;" | sudo -u postgres psql -U postgres postgres
sudo -u postgres pg_dump -t evidence > evidence.sql
sed -i -e 's/postgres/setflags/g' evidence.sql
source secrets/db.sh
echo "drop table evidence;" | psql
psql < evidence.sql
