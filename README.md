# `soda-rz`

<img src='graph.png'/>

Program to calculate and plot the probability distribution for
the total damage output (on the first round) of a [Soda Dungeon](http://sodadungeon.com)
party of 5 with 1 (or more) Ragezerker(s), taking into account the
stacked attack and probability to get critical hits, as the party
progresses through a given number of levels.

## Usage:

```
# install
$ git clone git@github.com:eugene-eeo/soda-rz.git
$ cd soda-rz

# make new virtual-env
$ pip install -r requirements.txt

# run simulation!
$ ./sim > data.txt
$ cat data.txt | ./gen graph.png --dpi=300
```

The party size, stats of party members (base damage, critical %,
and critical multiplier), etc can be easily configured by editing
the `config.json` file. A sample configuration is provided.
