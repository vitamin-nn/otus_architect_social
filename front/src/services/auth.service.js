import axios from 'axios';

const API_URL = process.env.VUE_APP_HTTP_SERVER_URL;

class AuthService {
    login(user) {
        return axios
            .post(API_URL + 'login', {
                email: user.email,
                password: user.password
            })
            .then(response => {
                if (response.data.access_token) {
                    localStorage.setItem('user', JSON.stringify(response.data.user));
                    localStorage.setItem('access_token', JSON.stringify(response.data.access_token));
                    localStorage.setItem('refresh_token', JSON.stringify(response.data.refresh_token));
                }

                return response.data;
            });
    }

    logout() {
        localStorage.removeItem('user');
        localStorage.removeItem('access_token');
        localStorage.removeItem('refresh_token');
    }

    register(user) {
        return axios
            .post(API_URL + 'register', {
                email: user.email,
                password: user.password,
                first_name: user.first_name,
                last_name: user.last_name,
                birth_date: user.birth_date,
                sex: user.sex,
                interest: user.interest,
                city: user.city
            })
            .then(response => {
                if (response.data.access_token) {
                    localStorage.setItem('user', JSON.stringify(response.data.user));
                    localStorage.setItem('access_token', JSON.stringify(response.data.access_token));
                    localStorage.setItem('refresh_token', JSON.stringify(response.data.refresh_token));
                }

                return response.data;
            });
    }
}

export default new AuthService();
