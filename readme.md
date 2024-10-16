# Fonctionnement de l'API

Multi-agent voting system with API implemented with Go, made in 2022.

## Installation du projet

Le projet est disponible sur le gitlab : https://gitlab.utc.fr/tleymone/API_Server.git

## Lancement

Pour lancer le serveur, il suffit de lancer la commande `go run bin/launch-server.go` dans un terminal.  
Ensuite, dans un autre terminal, vous pouvez lancer une démo avec `go run bin/demo/demo.go`.

## Fonctions implémentées

Une fois le serveur lancé, il y a 3 requêtes possibles.

### /new_ballot

Tout d'abord, pour créer un bulletin de vote, il faut utiliser une requête HTTP POST vers `/new_ballot` avec comme options :

- `rule` : la règle de vote parmi ['majority', 'borda', 'copeland', 'kemeny', 'stv', 'approval']
- `deadline` : la date de fin du vote sous forme de chaîne de caractères, de la forme RFC3339, c'est à dire "2012-11-01T22:08:41+01:00" par exemple
- `voter-ids` : le tableau contenant les ids des agents qui peuvent participer au vote, par exemple ["ag_id1", "ag_id2"]
- `#alts` : le nombre d'alternatives possible pour le vote

Si il n'y a pas d'erreurs, la requête renverra un objet JSON contenant `ballot-id`, l'ID du ballot créé.

### /vote

Ensuite, pour voter, il faut utiliser une requête HTTP POST vers `/vote` avec comme options :

- `agent-id` : l'ID du votant, cette ID doit être spécifié dans `voter-ids` à la création du bulletin de vote
- `vote-id` : l'ID du bulletin de vote, il est renvoyé par l'API à la création du bulletin de vote
- `prefs` : la liste de préférences du votant
- `options` : un tableau pouvant être vide qui sert à rajouter des paramètres si besoin, par exemple le seuil d'acceptation pour le vote approval

Si le le vote a bien été pris en compte, l'API renverra le code 200.

### /result

Pour finir, pour avoir le résultat du vote, il faut utiliser une requête HTTP POST vers `/result` avec comme option :

- `ballot-id` : l'ID du bulletin de vote

Le résultat ne peut se demander uniquement après la deadline fixé à la création du bulletin de vote sinon la requête renverra une erreur.

Si il n'y a pas d'erreurs, la requête renverra un objet JSON contenant :

- `winner` : l'alternative élue
- `ranking` : s'il existe, le classement des alternatives à l'issue du vote

## Fonctionnalités

L'ensemble des fonctionnalités de création de ballot, de vote, et d'obtention de résultat ont été implémentées dans ce projet.
Cependant, pour les procédure approval et kemeny, lors de l'obtention du résultat, le ranking n'est pas renvoyé. 
Le gagnant est tout de même renvoyé.

En cas d'égalité, la fonction de tiebreak renvoie le candidat avec l'indice le plus petit dans la liste des alternatives.
L'option pour changer la fonction de tiebreak n'a pas été implémentée.
