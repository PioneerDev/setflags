echo 'drop table if exists "user";' | sudo -u postgres psql -U postgres postgres
echo 'select user_id as id, identity_number, full_name, avatar_url, access_token, subscribed_at as created_at, subscribed_at as updated_at into table "user" from users u;' | sudo -u postgres psql -U postgres postgres
sudo -u postgres pg_dump -t user > user.sql
sed -i -e 's/postgres/setflags/g' user.sql
source secrets/db.sh
echo 'drop table if exists "user";' | psql
psql < user.sql
