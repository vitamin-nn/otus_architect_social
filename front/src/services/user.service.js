import axios from 'axios';
import authHeader from './auth-header';

const API_URL = process.env.VUE_APP_HTTP_SERVER_URL;

class UserService {
    getMainPageContent() {
        return axios.get(API_URL + '?limit=20&offset=0');
    }
    getPublicProfileById(id) {
        return axios.get(API_URL + 'user/' + id);
    }
    getMyProfile() {
        return axios.get(API_URL + 'profile', { headers: authHeader() });
    }
    addFriend(id) {
        return axios.post(API_URL + 'friends/add', { friend_id: id }, { headers: authHeader() });
    }
    removeFriend(id) {
        return axios.post(API_URL + 'friends/remove', { friend_id: id }, { headers: authHeader() });
    }
    getFriendList() {
        return axios.get(API_URL + 'friends?limit=20&offset=0', { headers: authHeader() });
    }
}

export default new UserService();
