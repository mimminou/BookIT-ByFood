module.exports = {
    async redirects() {
        return [
            {
                source: '/',
                destination: '/books',
                permanent: true
            }
        ]
    }
}
