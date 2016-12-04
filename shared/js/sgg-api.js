const superagent = require('superagent')

function requestWrap(request, version) {
    let r = request.set('Accept', `application/vnd.svgg.${version}+json`)

    let csrf = document.querySelector('meta[csrf]')
    if (csrf !== null) {
        r.set('CSRF-Token', csrf.getAttribute('csrf'))
    }

    return r.end()
}

module.exports = {

}
