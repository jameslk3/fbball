"use client";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogClose,
  DialogTitle,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogContent,
  DialogTrigger,
} from "@/components/ui/dialog";
import {
  Table,
  TableBody,
  TableCaption,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Eye, Trash2 } from "lucide-react";
import { Lineup, SlimGene, useLineup } from "@/app/context/LineupContext";
import { AnotherLineupCard } from "@/components/LineupDisplay";
import {
  Carousel,
  CarouselContent,
  CarouselItem,
  CarouselPrevious,
  CarouselNext,
} from "./ui/carousel";
import { ScrollArea } from "@/components/ui/scroll-area";

export function ManageLineupsTable() {
  const { savedLineups, deleteLineup } = useLineup();

  return (
    <Table className="w-full">
      <TableCaption>View and delete your saved lineups for the selected teams.</TableCaption>
      <TableHeader>
        <TableRow>
          <TableHead className="text-left w-[20%]">Timestamp</TableHead>
          <TableHead className="text-center w-[10%]">Improvement</TableHead>
          <TableHead className="text-center w-[10%]">Week</TableHead>
          <TableHead className="text-center w-[10%]">Threshold</TableHead>
          <TableHead className="text-center w-[10%]">View</TableHead>
          <TableHead className="text-center w-[10%]">Delete</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {savedLineups.length != 0 ? (
          <>
            {savedLineups.map((savedLineup) => (
              <TableRow key={savedLineup.Id}>
                <TableCell className="font-medium">
                  {savedLineup.Timestamp}
                </TableCell>
                <TableCell className="text-center">
                  {savedLineup.Improvement}
                </TableCell>
                <TableCell className="text-center">{savedLineup.Week}</TableCell>
                <TableCell className="text-center">{savedLineup.Threshold}</TableCell>
                <TableCell className="text-center">
                  <div className="inline-block">
                  <ViewLineupButton lineup={savedLineup} />
                  </div>
                </TableCell>
                <TableCell className="text-center">
                  <div className="inline-block">
                    <DeleteLineupConfirmation team_id={savedLineup.Id} />
                  </div>
                </TableCell>
              </TableRow>
            ))}
          </>
        ) : (
          <TableRow>
            <TableCell colSpan={5}>
              <p className="text-center text-gray-400">No lineups found.</p>
            </TableCell>
          </TableRow>
        )}
      </TableBody>
    </Table>
  );
}

function ViewLineupButton({ lineup }: { lineup: Lineup }) {
  return (
    <Dialog>
      <DialogTrigger>
        <Button variant="ghost" className="hover:bg-input">
          <Eye size={20} />
        </Button>
      </DialogTrigger>
      <DialogContent>
        <ScrollArea className="w-full h-[80vh]">
        <Carousel className="ml-3 w-[75%] max-w-xl mt-3">
          <CarouselContent>
            {lineup.Lineup.map((gene: SlimGene, index: number) => (
              <CarouselItem key={index}>
                <AnotherLineupCard gene={gene} />
              </CarouselItem>
            ))}
          </CarouselContent>
          <CarouselPrevious />
          <CarouselNext />
        </Carousel>
        </ScrollArea>
      </DialogContent>
    </Dialog>
  );
}


function DeleteLineupConfirmation({ team_id }: { team_id: number }) {
  const { deleteLineup } = useLineup();

  return (
    <Dialog>
      <DialogTrigger>
        <Button variant="ghost" className="hover:bg-input mr-[-5px]">
          <Trash2 size={20} />
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Delete Lineup</DialogTitle>
          <DialogDescription>
            Are you sure you want to delete this lineup?
          </DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <DialogClose asChild>
            <Button variant="outline" className="mr-2">
              Cancel
            </Button>
          </DialogClose>
          <Button onClick={() => deleteLineup(team_id)} variant="default">Delete</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
