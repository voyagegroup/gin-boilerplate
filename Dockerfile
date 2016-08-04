# ベースイメージは公式のDockerイメージを使っています。
# 各種設定のためのオプションについては以下を参照してください。
#
# https://hub.docker.com/_/mysql/
FROM mysql:5.6
EXPOSE 3306
CMD ["mysqld"]
