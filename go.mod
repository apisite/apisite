module github.com/apisite/apisite

go 1.12

replace (
	github.com/apisite/apitpl => ./../apitpl
	github.com/apisite/apitpl/ginapitpl => ./../apitpl/ginapitpl
	github.com/apisite/procapi => ./../procapi
	github.com/apisite/procapi/ginproc => ./../procapi/ginproc
	github.com/apisite/procapi/pgtype => ./../procapi/pgtype
)

require (
	github.com/acoshift/paginate v1.1.1
	github.com/apisite/apitpl v0.3.0
	github.com/apisite/apitpl/ginapitpl v0.3.0
	github.com/apisite/procapi v0.3.2
	github.com/apisite/procapi/ginproc v0.3.2
	github.com/apisite/procapi/pgtype v0.3.2
	github.com/birkirb/loggers-mapper-logrus v0.0.0-20180326232643-461f2d8e6f72
	github.com/gin-contrib/static v0.0.0-20190301062546-ed515893e96b
	github.com/gin-gonic/gin v1.3.0
	github.com/jessevdk/go-flags v1.4.0
	github.com/lib/pq v1.0.0
	github.com/pkg/errors v0.8.1
	github.com/sirupsen/logrus v1.4.1
	github.com/spf13/cast v1.3.0
	github.com/stretchr/testify v1.3.0
	gopkg.in/birkirb/loggers.v1 v1.1.0
)
