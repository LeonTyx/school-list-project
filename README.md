#School List Project
The School List Project aims to help parents and students by bringing together school supplies into one central location

#Create a postgresql database

#Generate Google oauth credentials

#Clone the repository
Use git to clone the repository
```git clone https://github.com/school-list-project```

#Create your environment.env file
Create a file called environment.env in the root of your project containing the following
```
GOOGLE_CLIENT_ID=<Google Client ID>
GOOGLE_CLIENT_SECRET=<Google Client Secret>
ENV=<DEV or PROD>
DATABASE_URL=postgres://<postgres>:<password>@<hostname>:<port>/<database>?sslmode=disable
DATABASE_SECRET=<Randomly generated string>
PORT=<8080>
```

#Compile the frontend
```npm run build```
#Compile the backend
```go build main.go```
