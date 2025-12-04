// src/app/deposits/page.js
"use client"

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger, DialogFooter } from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Badge } from "@/components/ui/badge";
import { Wallet, Plus, TrendingUp, ArrowUpRight } from "lucide-react";

export default function DepositsPage() {
  return (
    <div className="space-y-8 p-4 md:p-8 bg-[#F4F7FE] dark:bg-[#0B1437] min-h-screen">
      
      {/* HEADER & ADD BUTTON */}
      <div className="flex flex-col md:flex-row justify-between items-start md:items-center gap-4">
        <div>
           <h2 className="text-sm font-medium text-gray-500 dark:text-gray-400 mb-1">Overview</h2>
           <h1 className="text-3xl font-bold text-[#1B254B] dark:text-white tracking-tight">Deposits & Funds</h1>
        </div>

        {/* ADD DEPOSIT MODAL */}
        <Dialog>
          <DialogTrigger asChild>
            <Button className="bg-blue-600 hover:bg-blue-700 text-white rounded-xl px-6 h-12 shadow-lg shadow-blue-500/20">
              <Plus className="mr-2 h-4 w-4" /> Add New Deposit
            </Button>
          </DialogTrigger>
          <DialogContent className="sm:max-w-[425px] bg-white dark:bg-[#111C44] border-none rounded-2xl shadow-2xl">
            <DialogHeader>
              <DialogTitle className="text-[#1B254B] dark:text-white text-xl font-bold">Record New Deposit</DialogTitle>
            </DialogHeader>
            <div className="grid gap-6 py-4">
              <div className="grid gap-2">
                <Label htmlFor="name" className="text-gray-500">Member Name</Label>
                <Input id="name" placeholder="Search member..." className="rounded-xl bg-gray-50 dark:bg-[#0B1437] border-none h-12" />
              </div>
              <div className="grid gap-2">
                <Label htmlFor="amount" className="text-gray-500">Amount</Label>
                <div className="relative">
                    <span className="absolute left-3 top-3.5 text-gray-400">$</span>
                    <Input id="amount" placeholder="0.00" className="pl-7 rounded-xl bg-gray-50 dark:bg-[#0B1437] border-none h-12" />
                </div>
              </div>
              <div className="grid gap-2">
                <Label className="text-gray-500">Payment Method</Label>
                <Select>
                  <SelectTrigger className="rounded-xl bg-gray-50 dark:bg-[#0B1437] border-none h-12">
                    <SelectValue placeholder="Select method" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="cash">Cash</SelectItem>
                    <SelectItem value="bank">Bank Transfer</SelectItem>
                    <SelectItem value="check">Check</SelectItem>
                  </SelectContent>
                </Select>
              </div>
            </div>
            <DialogFooter>
              <Button type="submit" className="w-full bg-blue-600 hover:bg-blue-700 rounded-xl h-12 text-white font-bold">Confirm Deposit</Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>

      {/* KPI CARDS */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        
        {/* BIG CARD: Total Treasury */}
        <Card className="rounded-[20px] border-none shadow-sm bg-blue-600 dark:bg-blue-700 text-white overflow-hidden relative">
           <div className="absolute top-0 right-0 p-8 opacity-10">
              <Wallet className="h-48 w-48 text-white" />
           </div>
          <CardHeader className="pb-2">
            <CardTitle className="text-blue-100 font-medium text-sm">Total Treasury Balance</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-5xl font-bold tracking-tight mb-2">$8,450,200</div>
            <div className="flex items-center gap-2 text-blue-100">
               <span className="bg-white/20 px-2 py-1 rounded-lg text-xs font-bold flex items-center">
                  +3.5% <ArrowUpRight className="h-3 w-3 ml-1"/>
               </span>
               <span className="text-sm">vs last month</span>
            </div>
          </CardContent>
        </Card>

        {/* RECENT STATS */}
        <Card className="rounded-[20px] border-none shadow-sm bg-white dark:bg-[#111C44]">
          <CardHeader>
            <CardTitle className="text-gray-500 font-medium text-sm">This Month's Inflow</CardTitle>
          </CardHeader>
          <CardContent>
             <div className="flex items-center gap-4">
                <div className="h-14 w-14 bg-green-50 dark:bg-green-900/20 rounded-full flex items-center justify-center">
                    <TrendingUp className="h-7 w-7 text-green-600 dark:text-green-400" />
                </div>
                <div>
                    <div className="text-3xl font-bold text-[#1B254B] dark:text-white">$124,500</div>
                    <p className="text-gray-400 text-sm">32 Deposits recorded</p>
                </div>
             </div>
          </CardContent>
        </Card>
      </div>

      {/* DEPOSITS TABLE */}
      <Card className="rounded-[20px] border-none shadow-sm bg-white dark:bg-[#111C44]">
        <CardHeader>
            <CardTitle className="text-xl font-bold text-[#1B254B] dark:text-white">Recent Transactions</CardTitle>
        </CardHeader>
        <CardContent>
            <Table>
                <TableHeader>
                <TableRow className="hover:bg-transparent border-b border-gray-100 dark:border-gray-800">
                    <TableHead className="pl-4">Member</TableHead>
                    <TableHead>Date</TableHead>
                    <TableHead>Method</TableHead>
                    <TableHead>Status</TableHead>
                    <TableHead className="text-right pr-4">Amount</TableHead>
                </TableRow>
                </TableHeader>
                <TableBody>
                {/* Demo Row 1 */}
                <TableRow className="hover:bg-gray-50/50 dark:hover:bg-white/5 border-b border-gray-50 dark:border-gray-800 h-16">
                    <TableCell className="font-bold text-[#1B254B] dark:text-white pl-4">Sarah Jenkins</TableCell>
                    <TableCell className="text-gray-500">Today, 10:42 AM</TableCell>
                    <TableCell className="text-gray-500">Bank Transfer</TableCell>
                    <TableCell><Badge className="bg-green-100 text-green-700 dark:bg-green-900 dark:text-green-300 hover:bg-green-100 border-none">Completed</Badge></TableCell>
                    <TableCell className="text-right font-bold text-[#1B254B] dark:text-white pr-4">+$1,200.00</TableCell>
                </TableRow>
                {/* Demo Row 2 */}
                <TableRow className="hover:bg-gray-50/50 dark:hover:bg-white/5 border-none h-16">
                    <TableCell className="font-bold text-[#1B254B] dark:text-white pl-4">Michael Ross</TableCell>
                    <TableCell className="text-gray-500">Yesterday, 4:20 PM</TableCell>
                    <TableCell className="text-gray-500">Cash</TableCell>
                    <TableCell><Badge className="bg-green-100 text-green-700 dark:bg-green-900 dark:text-green-300 hover:bg-green-100 border-none">Completed</Badge></TableCell>
                    <TableCell className="text-right font-bold text-[#1B254B] dark:text-white pr-4">+$500.00</TableCell>
                </TableRow>
                </TableBody>
            </Table>
        </CardContent>
      </Card>
    </div>
  );
}