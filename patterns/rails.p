RUUID %{HEXDIGIT}{32}
# rails controller with action
RCONTROLLER (?P<controller>[^#]+)#(?P<action>\w+)

# this will often be the only line:
RAILS3HEAD (?m)Started %{WORD:verb} "%{URIPATHPARAM:request}" for %{IPORHOST:clientip} at (?P<timestamp>%{YEAR}-%{MONTHNUM}-%{MONTHDAY} %{HOUR}:%{MINUTE}:%{SECOND} %{ISO8601_TIMEZONE})
# for some a strange reason, params are stripped of {} - not sure that's a good idea.
RPROCESSING \W*Processing by %{RCONTROLLER} as (?P<format>\S+)(?:\W*Parameters: {%{DATA:params}}\W*)?
RAILS3FOOT Completed %{NUMBER:response}%{DATA} in %{NUMBER:totalms}ms %{RAILS3PROFILE}%{GREEDYDATA}
RAILS3PROFILE (?:\(Views: %{NUMBER:viewms}ms \| ActiveRecord: %{NUMBER:activerecordms}ms|\(ActiveRecord: %{NUMBER:activerecordms}ms)?

# putting it all together
RAILS3 %{RAILS3HEAD}(?:%{RPROCESSING})?(?P<context>(?:%{DATA}\n)*)(?:%{RAILS3FOOT})?
