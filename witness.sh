echo "drop table witness" | sudo -u postgres psql -U postgres postgres
# flag_id payee_id verified
echo "select p.packet_id as flag_id, p.user_id as payee_id, p.done as verified into table witness from participants p;" | sudo -u postgres psql -U postgres postgres
sudo -u postgres pg_dump -t witness > witness.sql
sed -i -e 's/postgres/setflags/g' witness.sql
source secrets/db.sh
echo "drop table witness;" | psql
psql < witness.sql
