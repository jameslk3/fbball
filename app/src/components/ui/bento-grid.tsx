"use client";
import { useEffect } from "react";
import { useLeague } from "@/components/LeagueContext";
import { cn } from "../../utils/cn";
import { WavyBackground } from "./wavy-background";
import { useTheme } from "next-themes";

export const BentoGrid = ({ className, children }: { className?: string; children?: React.ReactNode }) => {
  return (
    <div
      className={cn(
        "grid md:auto-rows-[18rem] grid-cols-1 md:grid-cols-3 gap-4 max-w-7xl mx-auto ",
        className
      )}
    >
      {children}
    </div>
  );
};

export const BentoGridItem = ({
  className,
  title,
  description,
  icon,
}: {
  className?: string;
  title?: string | React.ReactNode;
  description?: string | React.ReactNode;
  icon?: React.ReactNode;
}) => {

  const { theme } = useTheme();

  return (
    <WavyBackground
      containerClassName={cn(
        "row-span-1 rounded-xl group/bento hover:shadow-xl transition duration-200 shadow-input dark:shadow-none p-4 dark:bg-black dark:border-white/[0.2] bg-white border border-transparent justify-between flex flex-col space-y-0",
        className
      )}
      colors={["#c97826", "#faf7f9"]}
      waveOpacity={0.5}
      waveWidth={20}
      speed="slow"
      backgroundFill={theme === "dark" ? "black" : "white"}
    >
      <div className="group-hover/bento:translate-x-2 transition duration-200">
        {icon}
        <div className="font-sans font-bold text-neutral-600 dark:text-neutral-200 mb-2 mt-2">{title}</div>
        <div className="font-sans font-normal text-neutral-600 text-xs dark:text-neutral-300">{description}</div>
      </div>
    </WavyBackground>
  );
};