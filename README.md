
# ForumYnov

ForumYnov est une application web de forum développée en Go, conçue pour permettre aux utilisateurs de créer des sujets, de poster des messages et d'interagir au sein d'une communauté.

## 🚀 Fonctionnalités

- **Authentification utilisateur** : Inscription, connexion et déconnexion sécurisées.
- **Gestion des sujets** : Création, affichage et suppression de sujets de discussion.
- **Système de commentaires** : Les utilisateurs peuvent commenter les sujets existants.
- **Interface utilisateur responsive** : Conçue avec HTML, CSS et JavaScript pour une expérience utilisateur optimale.
- **Base de données SQLite** : Stockage des données utilisateur, des sujets et des commentaires.
- **Déploiement avec Docker** : Facilité de déploiement grâce à Docker et Docker Compose.

## 🛠️ Technologies utilisées

- **Langage principal** : Go (Golang)
- **Framework web** : net/http de Go
- **Base de données** : SQLite
- **Frontend** : HTML, CSS, JavaScript
- **Conteneurisation** : Docker, Docker Compose

## 📁 Structure du projet

```bash
ForumYnov/
├── database/           # Gestion de la base de données SQLite
├── handlers/           # Gestion des requêtes HTTP
├── models/             # Modèles de données
├── scripts/            # Scripts utilitaires
├── static/             # Fichiers statiques (CSS, JS, images)
├── templates/          # Templates HTML
├── tools/              # Outils supplémentaires
├── Dockerfile          # Fichier Docker pour l'image de l'application
├── docker-compose.yml  # Configuration Docker Compose
├── go.mod              # Fichier de gestion des dépendances Go
├── go.sum              # Sommes de contrôle des dépendances
└── main.go             # Point d'entrée de l'application
```

## 🐳 Déploiement avec Docker

1. **Cloner le dépôt :**

   ```bash
   git clone https://github.com/Ferevor/ForumYnov.git
   cd ForumYnov
   ```

2. **Construire et démarrer les conteneurs :**

   ```bash
   docker compose -f 'docker-compose.yml' up -d --build app
   ```

3. **Accéder à l'application :**

   Ouvrez votre navigateur et rendez-vous sur `http://localhost:8080` ou `http://127.0.0.1:8080` pour utiliser ForumYnov.