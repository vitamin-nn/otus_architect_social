export default function authHeader() {
    let accessToken = JSON.parse(localStorage.getItem('access_token'));

    if (accessToken) {
        return { Authorization: 'Bearer ' + accessToken };
    } else {
        return {};
    }
}
