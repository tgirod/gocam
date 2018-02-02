# TODO

- [ ] Reprendre l'algo d'import DXF pour en faire quelque chose de robuste
    - [ ] support arcs inversés
    - [ ] support splines
- [ ] Écrire un programme qui se contente de faire la conversion DXF -> gcode.

Pour plus tard :

- [ ] Algo offset simple
- [ ] Overcut (mouvement spécial)
- [ ] Algo offset avec les îles (but premier de ce programme)

# 2018.01.30

Pour construire des chemins à partir des entités DXF, on a besoin de savoir si
elles partagent une extrémité commune. Ca nécessite de scanner les extrémités
de toutes les entités pour reconstruire un graphe.

Dans certains cas (les arcs par exemple), l'entité DXF ne contient pas les
coordonnées des extrémités - celles-ci sont calculées lors de la conversion, et
ne sont donc pas strictement exactes. On ne peut donc pas se contenter d'un
test d'égalité. Dans ce cas, on cherche une extrémité suffisament proche au
sens d'une distance 2D.

Ca revient à chercher le point le plus proche dans un espace 2D. Pour ça, on
pourrait construire un Quadtree qui divise l'espace récursivement. Pour chaque
entité, on vérifie si il n'existe pas déjà une extrémité similaire. Si oui, on
utilise celle-là pour représenter l'entité dans le modèle. Sinon, on ajoute une
nouvelle extrémité au Quadtree.

# 2018.01.24

Je me heurte à un problème pour faire de la gravure : les fichiers vectorisés
de Tim contiennent des splines, qui ne sont pas supportés par bCNC et cammill.
Du coup, la procédure est la suivante :

1. convertir le SVG en DXF avec inkscape,
2. ouvrir le DXF avec librecad et «exploser» le modèle, pour remplacer toutes
   les courbes par des segments,
3. convertir le DXF explosé en Gcode.

Cette procédure a cependant plusieurs problèmes :

1. l'explosion génère des fichiers trop lourds (librecad ne permet pas de
   contrôler la finesse de l'explosion),
2. bCNC et cammill sont incapables de tenir la charge avec des fichiers aussi
   lourds.

C'est là que gocam entre en jeu. J'ai essayé de l'utiliser pour faire la
conversion DXF -> Gcode, et c'est pas loin de marcher ! Il y a encore un peu de
boulot, mais au moins mon implémentation sort quelque chose sans faire planter
ma machine. 

Conclusion : Reprendre le code dxf -> gcode, et en faire un outil robuste -
c'est déjà une fonctionnalité utile.
