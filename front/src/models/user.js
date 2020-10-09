export default class User {
    constructor(email, password, first_name, last_name, birth_date, sex, interest, city) {
        this.email = email;
        this.password = password;
        this.first_name = first_name;
        this.last_name = last_name;
        this.birth_date = birth_date;
        this.sex = sex;
        this.interest = interest;
        this.city = city;
    }
}