module github.com/apisite/apisite

go 1.12

replace (
	github.com/apisite/pgcall => ../pgcall
	github.com/apisite/pgcall/gin-pgcall => ../pgcall/gin-pgcall
	github.com/apisite/pgcall/pgx-pgcall => ../pgcall/pgx-pgcall
	github.com/apisite/tpl2x => ../tpl2x
	github.com/apisite/tpl2x/gin-tpl2x => ../tpl2x/gin-tpl2x
)

require (
	github.com/acoshift/paginate v1.1.1
	github.com/apisite/pgcall v0.0.0-00010101000000-000000000000
	github.com/apisite/pgcall/gin-pgcall v0.0.0
	github.com/apisite/pgcall/pgx-pgcall v0.0.0
	github.com/apisite/tpl2x v0.0.0-20190323161051-eda7d2ca63fb
	github.com/apisite/tpl2x/gin-tpl2x v0.0.0-00010101000000-000000000000
	github.com/birkirb/loggers-mapper-logrus v0.0.0-20180326232643-461f2d8e6f72
	github.com/gin-contrib/static v0.0.0-20190301062546-ed515893e96b
	github.com/gin-gonic/gin v1.3.0
	github.com/jessevdk/go-flags v1.4.0
	github.com/sirupsen/logrus v1.4.0
	github.com/spf13/cast v1.3.0
)
