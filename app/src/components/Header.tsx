"use client";
import Image from "next/image";
import { ModeToggle } from "./ui/toggle-mode";

export default function Header() {
  return (
    <header>
      <nav>
        <ul className="flex items-center justify-between">
          <li className="flex-shrink-0.5">
						<Image
							src="/logo.png"
							alt="Logo"
							width={120}
							height={100}
							priority
						/>
          </li>
          <hr className="w-1/4 border-primary flex-shrink-1 flex-grow-1"></hr>
          {/* font-bold pb-1 border-b-4 border-primary */}
					<li className="text-4xl font-bold flex-shrink-0">
            Court Visionaries
          </li>
          <hr className="w-1/4 border-primary flex-shrink-1 flex-grow-1"></hr>
					<li className="flex-shrink-0 mr-20"><ModeToggle /></li>
        </ul>
      </nav>
    </header>
  );
}