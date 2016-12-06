const sgg = require('./lib/sgg-api')

sgg.getUser({slug: 'kts', 'other': 'yep'}).then((data) => {
    console.log(data)
}).catch((err) => {
    console.error('[ERR]', err)
})
