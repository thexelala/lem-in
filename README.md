# Lem-in

*Lem-in* is a project that focuses on efficiently solving a movement problem within a network. The problem involves moving a group of ants from a start room to an end room through a network of rooms and tunnels, while minimizing the time required for the movement. Each tunnel and room can accommodate at most one ant at a time.
To achieve optimal movement, we first find all possible paths through the network. Then, we use specific criteria, such as the overlap rate and path length, to select the most efficient paths while ensuring they adhere to the defined constraints.


## Features

- **Parsing:** The program must effectively parse the input describing the network of rooms and tunnels.
- **Optimal Path Search:** The algorithm must be able to find the most efficient set of paths to move the ants from the start room to the end room.
- **Ants Management:** The program must manage the movement of ants through the network of rooms, ensuring that no room is overloaded and optimizing their path.
- **Movement Optimization:** The goal is to minimize the time required for movement by judiciously distributing the ants over the different available paths.
- **Ants Movement Display:** Display the movements of each ant on standard output. The format is `Lx-y` where `x` is the ant's name and `y` is its destination room.

## Installation

```bash
git clone https://zone01normandie.org/git/nduval/lem-in
```

## Usage

```bash
cd lem-in/cmd
go run . [file_name]
```